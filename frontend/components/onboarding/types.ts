// Onboarding Types & Step Configuration

export interface OnboardingPayload {
  name: string;
  allergens: string[];
  dietary_preferences: string[];
  nutritional_goals: string[];
  cooking_skill: string;
  cuisine_preferences: string[];
}

export type StepType = 'text' | 'bubble-multi' | 'stacked-single';

export interface OnboardingStep {
  id: number;
  label: string;
  title: string;
  subtitle?: string;
  field: keyof OnboardingPayload;
  type: StepType;
  options?: string[];
}

export const ONBOARDING_STEPS: OnboardingStep[] = [
  {
    id: 1,
    label: 'IDENTIFY',
    title: 'What should we call you?',
    subtitle: 'Your name helps us personalize your experience.',
    field: 'name',
    type: 'text',
  },
  {
    id: 2,
    label: 'ALLERGIES',
    title: 'Any food allergies?',
    subtitle: "We'll make sure to flag these in your recommendations.",
    field: 'allergens',
    type: 'bubble-multi',
    options: [
      'Dairy',
      'Eggs',
      'Peanuts',
      'Tree Nuts',
      'Shellfish',
      'Fish',
      'Wheat',
      'Soy',
      'Sesame',
    ],
  },
  {
    id: 3,
    label: 'DIET',
    title: 'Dietary preferences?',
    subtitle: 'Select all that apply to your lifestyle.',
    field: 'dietary_preferences',
    type: 'bubble-multi',
    options: [
      'Vegetarian',
      'Vegan',
      'Pescatarian',
      'Keto',
      'Paleo',
      'Halal',
      'Kosher',
      'Gluten-Free',
    ],
  },
  {
    id: 4,
    label: 'GOALS',
    title: 'Nutritional goals?',
    subtitle: 'What are you optimizing for?',
    field: 'nutritional_goals',
    type: 'bubble-multi',
    options: [
      'High-Protein',
      'Low-Carb',
      'Low-Calorie',
      'Low-Sodium',
      'Heart-Healthy',
      'Balanced',
    ],
  },
  {
    id: 5,
    label: 'SKILL',
    title: 'Cooking experience?',
    subtitle: "We'll tailor recipe complexity to your level.",
    field: 'cooking_skill',
    type: 'stacked-single',
    options: ['Beginner', 'Intermediate', 'Advanced'],
  },
  {
    id: 6,
    label: 'CUISINES',
    title: 'Favorite cuisines?',
    subtitle: 'Select all that excite your palate.',
    field: 'cuisine_preferences',
    type: 'bubble-multi',
    options: [
      'Italian',
      'Mexican',
      'Chinese',
      'Japanese',
      'Indian',
      'Thai',
      'Mediterranean',
      'French',
      'American',
      'Korean',
      'Vietnamese',
      'Middle Eastern',
    ],
  },
];

export const initialPayload: OnboardingPayload = {
  name: '',
  allergens: [],
  dietary_preferences: [],
  nutritional_goals: [],
  cooking_skill: '',
  cuisine_preferences: [],
};
