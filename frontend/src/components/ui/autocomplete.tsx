import { useState, useEffect, useRef } from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useDebounce } from "@/hooks/useDebounce";

interface AutocompleteItem {
  id: string;
  name: string;
}

interface AutocompleteProps {
  value?: string;
  onValueChange: (value: string) => void;
  searchFn: (query: string) => Promise<AutocompleteItem[]>;
  placeholder?: string;
  emptyMessage?: string;
  disabled?: boolean;
}

export function Autocomplete({
  value = "",
  onValueChange,
  searchFn,
  placeholder = "Search...",
  emptyMessage = "No results found.",
  disabled = false,
}: AutocompleteProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState(value);
  const [results, setResults] = useState<AutocompleteItem[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const debouncedQuery = useDebounce(query, 300);
  const abortControllerRef = useRef<AbortController | null>(null);

  useEffect(() => {
    const search = async () => {
      if (debouncedQuery.length < 2) {
        setResults([]);
        return;
      }

      // Cancel previous request
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }

      abortControllerRef.current = new AbortController();
      setIsLoading(true);

      try {
        const items = await searchFn(debouncedQuery);
        setResults(items);
      } catch (error: any) {
        if (error.name !== "AbortError") {
          console.error("Search failed:", error);
          setResults([]);
        }
      } finally {
        setIsLoading(false);
      }
    };

    search();

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [debouncedQuery, searchFn]);

  const handleSelect = (item: AutocompleteItem) => {
    setQuery(item.name);
    onValueChange(item.name);
    setOpen(false);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          disabled={disabled}
          className="w-full justify-between"
        >
          {query || placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0" align="start">
        <Command shouldFilter={false}>
          <CommandInput
            placeholder={placeholder}
            value={query}
            onValueChange={setQuery}
          />
          <CommandList>
            {isLoading && (
              <div className="py-6 text-center text-sm">Searching...</div>
            )}
            {!isLoading && query.length < 2 && (
              <div className="py-6 text-center text-sm text-muted-foreground">
                Type at least 2 characters to search
              </div>
            )}
            {!isLoading && query.length >= 2 && results.length === 0 && (
              <CommandEmpty>{emptyMessage}</CommandEmpty>
            )}
            {!isLoading && results.length > 0 && (
              <CommandGroup>
                {results.map((item) => (
                  <CommandItem
                    key={item.id}
                    value={item.name}
                    onSelect={() => handleSelect(item)}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        query === item.name ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {item.name}
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
