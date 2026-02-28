'use client';

import { useState, useCallback, useEffect, Suspense } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import {
  NameInput,
  BubbleSelect,
  StackedSelect,
  ONBOARDING_STEPS,
  initialPayload,
  type OnboardingPayload,
  type OnboardingStep,
} from '@/components/onboarding';
import { getUserByAuth0ID } from '@/lib/user-api';
import { useUser } from '@/contexts/user-context';

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

// Animation configuration
const slideVariants = {
  enter: (direction: number) => ({
    x: direction > 0 ? 80 : -80,
    opacity: 0,
  }),
  center: {
    x: 0,
    opacity: 1,
  },
  exit: (direction: number) => ({
    x: direction < 0 ? 80 : -80,
    opacity: 0,
  }),
};

const slideTransition = {
  x: { type: 'tween', duration: 0.4, ease: [0.22, 1, 0.36, 1] },
  opacity: { duration: 0.3 },
};

function OnboardingContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const uid = searchParams.get('uid');
  const { user: contextUser } = useUser();

  const [currentStep, setCurrentStep] = useState(0);
  const [direction, setDirection] = useState(0);
  const [payload, setPayload] = useState<OnboardingPayload>(initialPayload);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Resolve auth0 sub → internal UUID
  const [internalUserId, setInternalUserId] = useState<string | null>(null);

  useEffect(() => {
    // Prefer the user already in context (redirected from UserProvider)
    if (contextUser?.id) {
      setInternalUserId(contextUser.id);
      return;
    }
    // Fall back to ?uid= auth0 sub param (fresh OAuth flow)
    if (!uid) return;
    getUserByAuth0ID(uid)
      .then((user) => {
        setInternalUserId(user.id);
        // Also cache in localStorage so the user context picks it up after redirect
        try {
          localStorage.setItem('sift_user', JSON.stringify(user));
        } catch {
          // ignore
        }
      })
      .catch((err) => {
        console.error('Failed to resolve user:', err);
        setError('Could not identify your account. Please try logging in again.');
      });
  }, [uid, contextUser?.id]);

  const step = ONBOARDING_STEPS[currentStep];
  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === ONBOARDING_STEPS.length - 1;
  const totalSteps = ONBOARDING_STEPS.length;

  // Update payload for a specific field
  const updateField = useCallback(
    <K extends keyof OnboardingPayload>(field: K, value: OnboardingPayload[K]) => {
      setPayload((prev) => ({ ...prev, [field]: value }));
    },
    []
  );

  // Submit onboarding data
  const handleSubmit = useCallback(async () => {
    if (!internalUserId) {
      setError('Missing user ID. Please try logging in again.');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    // Filter out empty fields
    const filteredPayload: Partial<OnboardingPayload> = {};
    if (payload.name.trim()) filteredPayload.name = payload.name.trim();
    if (payload.allergens.length > 0) filteredPayload.allergens = payload.allergens;
    if (payload.dietary_preferences.length > 0)
      filteredPayload.dietary_preferences = payload.dietary_preferences;
    if (payload.nutritional_goals.length > 0)
      filteredPayload.nutritional_goals = payload.nutritional_goals;
    if (payload.cooking_skill) filteredPayload.cooking_skill = payload.cooking_skill.toLowerCase();
    if (payload.cuisine_preferences.length > 0)
      filteredPayload.cuisine_preferences = payload.cuisine_preferences;

    try {
      const response = await fetch(`${API_URL}/users/${internalUserId}/onboarding`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(filteredPayload),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `Failed to save preferences (${response.status})`);
      }

      // Update localStorage with the fresh user from the backend (onboarding_completed: true)
      try {
        const updatedUser = await response.json();
        if (updatedUser?.id) {
          localStorage.setItem('sift_user', JSON.stringify(updatedUser));
        }
      } catch {
        // ignore if response has no body
      }

      // Success - redirect to dashboard
      router.push('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unexpected error occurred');
      setIsSubmitting(false);
    }
  }, [internalUserId, payload, router]);

  // Navigate to next step
  const handleNext = useCallback(async () => {
    if (isLastStep) {
      await handleSubmit();
    } else {
      setDirection(1);
      setCurrentStep((prev) => prev + 1);
    }
  }, [isLastStep, handleSubmit]);

  // Navigate to previous step
  const handleBack = useCallback(() => {
    if (!isFirstStep) {
      setDirection(-1);
      setCurrentStep((prev) => prev - 1);
    }
  }, [isFirstStep]);

  // Render the current step content
  const renderStepContent = (step: OnboardingStep) => {
    switch (step.type) {
      case 'text':
        return (
          <NameInput
            value={payload.name}
            onChange={(value) => updateField('name', value)}
          />
        );
      case 'bubble-multi':
        return (
          <BubbleSelect
            options={step.options || []}
            selected={payload[step.field] as string[]}
            onChange={(value) => updateField(step.field, value as never)}
          />
        );
      case 'stacked-single':
        return (
          <StackedSelect
            options={step.options || []}
            selected={payload[step.field] as string}
            onChange={(value) => updateField(step.field, value as never)}
          />
        );
      default:
        return null;
    }
  };

  // Handle Enter key on name input
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !isSubmitting) {
      handleNext();
    }
  };

  return (
    <div className="min-h-screen bg-[#F9F8F6] flex items-center justify-center p-4 md:p-8">
      {/* The Ledger Container */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
        className="w-full max-w-2xl border border-[#2A2724]/20 bg-[#F9F8F6]"
        onKeyDown={handleKeyDown}
      >
        {/* Header - Progress Indicator */}
        <div className="px-8 py-6 border-b border-[#2A2724]/10">
          <div className="flex items-center justify-between">
            <span className="font-sans text-xs tracking-[0.2em] uppercase text-[#2A2724]/50">
              Step {String(step.id).padStart(2, '0')} — {step.label}
            </span>
            <span className="font-sans text-xs tracking-[0.2em] uppercase text-[#2A2724]/30">
              {currentStep + 1} / {totalSteps}
            </span>
          </div>
          
          {/* Progress bar */}
          <div className="mt-4 h-px bg-[#2A2724]/10 relative overflow-hidden">
            <motion.div
              className="absolute inset-y-0 left-0 bg-[#4A5D4E]"
              initial={{ width: 0 }}
              animate={{ width: `${((currentStep + 1) / totalSteps) * 100}%` }}
              transition={{ duration: 0.4, ease: [0.22, 1, 0.36, 1] }}
            />
          </div>
        </div>

        {/* Dynamic Content Area */}
        <div className="px-8 py-12 min-h-[400px] flex flex-col">
          {/* Question Title */}
          <AnimatePresence mode="wait" custom={direction}>
            <motion.div
              key={`title-${currentStep}`}
              custom={direction}
              variants={slideVariants}
              initial="enter"
              animate="center"
              exit="exit"
              transition={slideTransition}
              className="mb-8"
            >
              <h1 className="font-serif text-3xl md:text-4xl text-[#2A2724] tracking-[-0.02em] leading-tight">
                {step.title}
              </h1>
              {step.subtitle && (
                <p className="mt-3 font-sans text-sm text-[#2A2724]/50 tracking-[0.01em]">
                  {step.subtitle}
                </p>
              )}
            </motion.div>
          </AnimatePresence>

          {/* Step Content */}
          <div className="flex-1">
            <AnimatePresence mode="wait" custom={direction}>
              <motion.div
                key={`content-${currentStep}`}
                custom={direction}
                variants={slideVariants}
                initial="enter"
                animate="center"
                exit="exit"
                transition={slideTransition}
              >
                {renderStepContent(step)}
              </motion.div>
            </AnimatePresence>
          </div>

          {/* Error Message */}
          <AnimatePresence>
            {error && (
              <motion.p
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="mt-6 font-sans text-sm text-red-600"
              >
                {error}
              </motion.p>
            )}
          </AnimatePresence>
        </div>

        {/* Footer - Navigation */}
        <div className="px-8 py-6 border-t border-[#2A2724]/10 flex items-center justify-between">
          {/* Back Button */}
          <button
            type="button"
            onClick={handleBack}
            disabled={isFirstStep || isSubmitting}
            className={`
              font-sans text-xs tracking-[0.2em] uppercase
              transition-all duration-200
              ${
                isFirstStep || isSubmitting
                  ? 'text-[#2A2724]/20 cursor-not-allowed'
                  : 'text-[#2A2724]/60 hover:text-[#2A2724] cursor-pointer'
              }
            `}
          >
            Back
          </button>

          {/* Skip indicator */}
          {!isLastStep && (
            <span className="font-sans text-[10px] tracking-[0.15em] uppercase text-[#2A2724]/30">
              Optional — skip anytime
            </span>
          )}

          {/* Next / Complete Button */}
          <motion.button
            type="button"
            onClick={handleNext}
            disabled={isSubmitting}
            whileHover={{ scale: isSubmitting ? 1 : 1.02 }}
            whileTap={{ scale: isSubmitting ? 1 : 0.98 }}
            className={`
              px-6 py-3
              font-sans text-xs tracking-[0.2em] uppercase
              border border-[#2A2724]/30
              transition-all duration-200
              ${
                isSubmitting
                  ? 'text-[#2A2724]/40 cursor-wait'
                  : 'text-[#2A2724] hover:bg-[#2A2724] hover:text-[#F9F8F6] cursor-pointer'
              }
            `}
          >
            {isSubmitting ? (
              <span className="flex items-center gap-2">
                <motion.span
                  animate={{ rotate: 360 }}
                  transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                  className="inline-block w-3 h-3 border border-current border-t-transparent rounded-full"
                />
                Saving...
              </span>
            ) : isLastStep ? (
              'Complete Profile'
            ) : (
              'Next'
            )}
          </motion.button>
        </div>
      </motion.div>
    </div>
  );
}

// Loading fallback
function OnboardingLoading() {
  return (
    <div className="min-h-screen bg-[#F9F8F6] flex items-center justify-center">
      <div className="font-sans text-xs tracking-[0.2em] uppercase text-[#2A2724]/40">
        Loading...
      </div>
    </div>
  );
}

// Main page component with Suspense boundary
export default function OnboardingPage() {
  return (
    <Suspense fallback={<OnboardingLoading />}>
      <OnboardingContent />
    </Suspense>
  );
}
