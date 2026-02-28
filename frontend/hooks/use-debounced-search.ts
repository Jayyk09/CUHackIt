'use client'

import { useState, useEffect, useRef } from 'react'
import { FoodItem, mockSearchProducts } from '@/lib/mock-data'

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
      // In the future, this will be an API call:
      // const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`)
      // const data = await response.json()
      // setResults(data.results)
      
      // For now, use mock data
      const searchResults = mockSearchProducts(query)
      setResults(searchResults)
      setIsLoading(false)
    }, 300)

    // Cleanup
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [query, isSearching])

  return { results, isLoading, isSearching }
}
