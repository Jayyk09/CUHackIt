'use client'

import { useState, useEffect, useRef } from 'react'
import { FoodItem, searchFoodProducts } from '@/lib/food-api'

interface UseDebouncedSearchResult {
  results: FoodItem[]
  isLoading: boolean
  isSearching: boolean // true when query >= 3 chars
}

/**
 * Debounced search hook
 * - Waits 300ms after last keystroke before searching
 * - Only searches when query is 3+ characters
 * - Returns loading state and results
 */
export function useDebouncedSearch(query: string): UseDebouncedSearchResult {
  const [results, setResults] = useState<FoodItem[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)
  const abortRef = useRef<AbortController | null>(null)

  const isSearching = query.length >= 3

  useEffect(() => {
    // Clear any existing timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }

    // If query is less than 3 characters, clear results
    if (!isSearching) {
      setResults([])
      setIsLoading(false)
      return
    }

    // Set loading state
    setIsLoading(true)

    // Debounce the search
    timeoutRef.current = setTimeout(() => {
      abortRef.current?.abort()
      const controller = new AbortController()
      abortRef.current = controller

      console.info('search request start', { query })
      searchFoodProducts({
        search: query,
        limit: 12,
        offset: 0,
        signal: controller.signal,
      })
        .then((searchResults) => {
          console.info('search request success', { query, count: searchResults.length })
          setResults(searchResults)
        })
        .catch((err: unknown) => {
          if (err instanceof Error && err.name === 'AbortError') {
            return
          }
          console.error('search request failed', { query, err })
          setResults([])
        })
        .finally(() => {
          setIsLoading(false)
        })
    }, 300)

    // Cleanup
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
      abortRef.current?.abort()
    }
  }, [query, isSearching])

  return { results, isLoading, isSearching }
}
