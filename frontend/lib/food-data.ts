export interface FoodItem {
  id: number
  name: string
  label: string
  src: string
  category: string
}

export const foodItems: FoodItem[] = [
  // Protein
  {
    id: 1,
    name: 'canned-tuna',
    label: 'Canned Tuna',
    src: '/realfood/001_canned-tuna.webp.png',
    category: 'Protein',
  },
  {
    id: 2,
    name: 'eggs',
    label: 'Eggs',
    src: '/realfood/002_eggs.webp.png',
    category: 'Protein',
  },
  {
    id: 4,
    name: 'ground-beef',
    label: 'Ground Beef',
    src: '/realfood/004_ground-beef.webp.png',
    category: 'Protein',
  },
  {
    id: 7,
    name: 'salmon',
    label: 'Salmon',
    src: '/realfood/007_salmon.webp.png',
    category: 'Protein',
  },
  {
    id: 16,
    name: 'steak',
    label: 'Steak',
    src: '/realfood/016_steak.webp.png',
    category: 'Protein',
  },
  {
    id: 43,
    name: 'chicken',
    label: 'Chicken',
    src: '/realfood/043_chicken.webp.png',
    category: 'Protein',
  },
  {
    id: 53,
    name: 'shrimp',
    label: 'Shrimp',
    src: '/realfood/053_shrimp.webp.png',
    category: 'Protein',
  },

  // Dairy & Fats
  {
    id: 6,
    name: 'cheese',
    label: 'Cheese',
    src: '/realfood/006_cheese.webp.png',
    category: 'Dairy',
  },
  {
    id: 9,
    name: 'yogurt',
    label: 'Yogurt',
    src: '/realfood/009_yogurt.webp.png',
    category: 'Dairy',
  },
  {
    id: 21,
    name: 'butter',
    label: 'Butter',
    src: '/realfood/021_butter.webp.png',
    category: 'Dairy',
  },
  {
    id: 30,
    name: 'milk',
    label: 'Milk',
    src: '/realfood/030_milk.webp.png',
    category: 'Dairy',
  },
  {
    id: 36,
    name: 'olive-oil',
    label: 'Olive Oil',
    src: '/realfood/036_olive-oil.webp.png',
    category: 'Fats',
  },
  {
    id: 5,
    name: 'avocado',
    label: 'Avocado',
    src: '/realfood/005_avocado.webp.png',
    category: 'Fats',
  },

  // Nuts & Seeds
  {
    id: 12,
    name: 'walnut',
    label: 'Walnut',
    src: '/realfood/012_walnut-shelled.webp.png',
    category: 'Nuts',
  },
  {
    id: 14,
    name: 'peanuts',
    label: 'Peanuts',
    src: '/realfood/014_peanuts.webp.png',
    category: 'Nuts',
  },
  {
    id: 52,
    name: 'almond',
    label: 'Almond',
    src: '/realfood/052_almond.webp.png',
    category: 'Nuts',
  },

  // Vegetables
  {
    id: 8,
    name: 'broccoli',
    label: 'Broccoli',
    src: '/realfood/008_broccoli.webp.png',
    category: 'Vegetable',
  },
  {
    id: 10,
    name: 'tomatoes',
    label: 'Tomatoes',
    src: '/realfood/010_tomatoes.webp.png',
    category: 'Vegetable',
  },
  {
    id: 24,
    name: 'carrots',
    label: 'Carrots',
    src: '/realfood/024_carrots.webp.png',
    category: 'Vegetable',
  },
  {
    id: 32,
    name: 'potato',
    label: 'Potato',
    src: '/realfood/032_potato.webp.png',
    category: 'Vegetable',
  },
  {
    id: 56,
    name: 'green-beans',
    label: 'Green Beans',
    src: '/realfood/056_green-beans.webp.png',
    category: 'Vegetable',
  },
  {
    id: 77,
    name: 'lettuce',
    label: 'Lettuce',
    src: '/realfood/077_lettuce.webp.png',
    category: 'Vegetable',
  },
  {
    id: 15,
    name: 'butternut',
    label: 'Butternut Squash',
    src: '/realfood/015_butternut.webp.png',
    category: 'Vegetable',
  },
  {
    id: 28,
    name: 'frozen-peas',
    label: 'Frozen Peas',
    src: '/realfood/028_frozen-peas.webp.png',
    category: 'Vegetable',
  },

  // Fruits
  {
    id: 3,
    name: 'blueberry',
    label: 'Blueberries',
    src: '/realfood/003_blueberry.webp.png',
    category: 'Fruit',
  },
  {
    id: 11,
    name: 'blueberry-alt',
    label: 'Blueberry',
    src: '/realfood/011_blueberry.webp.png',
    category: 'Fruit',
  },
  {
    id: 13,
    name: 'strawberry',
    label: 'Strawberry',
    src: '/realfood/013_strawberry-right.webp.png',
    category: 'Fruit',
  },
  {
    id: 44,
    name: 'grapes',
    label: 'Grapes',
    src: '/realfood/044_grapes.webp.png',
    category: 'Fruit',
  },
  {
    id: 45,
    name: 'bananas',
    label: 'Bananas',
    src: '/realfood/045_bananas.webp.png',
    category: 'Fruit',
  },
  {
    id: 60,
    name: 'apples',
    label: 'Apples',
    src: '/realfood/060_apples.webp.png',
    category: 'Fruit',
  },
  {
    id: 62,
    name: 'oranges',
    label: 'Oranges',
    src: '/realfood/062_oranges.webp.png',
    category: 'Fruit',
  },

  // Whole Grains
  {
    id: 23,
    name: 'bowl-oats',
    label: 'Oatmeal',
    src: '/realfood/023_bowl-oats.webp.png',
    category: 'Grain',
  },
  {
    id: 49,
    name: 'rice-beans',
    label: 'Rice & Beans',
    src: '/realfood/049_bowl-rice-beans.webp.png',
    category: 'Grain',
  },
  {
    id: 63,
    name: 'bread',
    label: 'Bread',
    src: '/realfood/063_bread.webp.png',
    category: 'Grain',
  },
  {
    id: 84,
    name: 'oats',
    label: 'Oats',
    src: '/realfood/084_oats.webp.png',
    category: 'Grain',
  },
]

export const foodByCategory = {
  protein: foodItems.filter(f => f.category === 'Protein'),
  dairy: foodItems.filter(f => f.category === 'Dairy'),
  fats: foodItems.filter(f => f.category === 'Fats'),
  nuts: foodItems.filter(f => f.category === 'Nuts'),
  vegetable: foodItems.filter(f => f.category === 'Vegetable'),
  fruit: foodItems.filter(f => f.category === 'Fruit'),
  grain: foodItems.filter(f => f.category === 'Grain'),
}

export const getFoodByName = (name: string) => foodItems.find(f => f.name === name)
export const getFoodById = (id: number) => foodItems.find(f => f.id === id)
