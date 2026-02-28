export function SiteFooter() {
  return (
    <footer className="border-t border-foreground px-4 md:px-8 lg:px-16 py-8 md:py-12">
      <div className="flex flex-col md:flex-row items-start md:items-end justify-between gap-6">
        <div>
          <span className="font-serif text-2xl md:text-3xl tracking-[-0.02em] text-foreground block">
            Sift
          </span>
          <p className="mt-2 font-sans text-xs tracking-[0.1em] uppercase text-muted-foreground">
            Real food data for real people.
          </p>
        </div>
        <div className="flex items-center gap-8">
          <a
            href="#"
            className="font-sans text-[10px] tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
          >
            Privacy
          </a>
          <a
            href="#"
            className="font-sans text-[10px] tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
          >
            Terms
          </a>
          <a
            href="#"
            className="font-sans text-[10px] tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
          >
            Source
          </a>
        </div>
      </div>
      <div className="mt-8 pt-6 border-t border-foreground/10">
        <p className="font-sans text-[10px] tracking-[0.15em] uppercase text-muted-foreground">
          Data sourced from USDA FoodData Central. Not medical advice.
        </p>
      </div>
    </footer>
  )
}
