import { api } from "@/lib/api";
import { sanitizeSearchQuery } from "@/utils/validation";
import { type AutocompleteItem } from "./master.data.service";

// ==================== FRONTEND INTERFACES ====================

export interface Skill {
  id: string;
  name: string;
}

// ==================== SERVICES ====================

export const skillService = {
  /**
   * Search skills from backend (case-insensitive)
   * Returns AutocompleteItem[] for compatibility with Autocomplete/Command components
   * @param query - Search term
   */
  search: async (query: string): Promise<AutocompleteItem[]> => {
    if (!query || query.length < 1) {
      return [];
    }

    try {
      const sanitized = sanitizeSearchQuery(query);
      // Backend returns array directly: [{ID: string, Name: string}]
      const response = await api.get<{ ID: string; Name: string }[]>("/skill", {
        params: { q: sanitized },
      });

      // Map backend response (PascalCase) to AutocompleteItem (camelCase)
      return (response.data || []).map((item) => ({
        id: item.ID,
        name: item.Name,
      }));
    } catch (error) {
      console.error("Failed to search skills:", error);
      return [];
    }
  },

  /**
   * Create a new skill (Optional: if you want to support explicit creation via API)
   * Note: The current backend design might handle creation implicitly or via a separate endpoint
   */
  create: async (name: string): Promise<Skill | null> => {
    try {
      // Assuming POST /skill creates a new skill
      const response = await api.post<Skill>("/skill", { name });
      return response.data;
    } catch (error) {
      console.error("Failed to create skill:", error);
      return null;
    }
  },
};
