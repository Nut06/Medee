import { api } from "@/lib/api";
import { sanitizeSearchQuery } from "@/utils/validation";

// ==================== AUTOCOMPLETE INTERFACE ====================

// Interface expected by the Autocomplete component
export interface AutocompleteItem {
  id: string;
  name: string;
}

// ==================== FRONTEND INTERFACES ====================

export interface Institute {
  id?: string;
  name: string;
  country?: string;
  domain?: string;
  webPages?: string[];
  source?: string;
}

export interface FieldOfStudy {
  id?: string;
  code: string;
  title: string;
  name?: string; // Legacy field from old database
  category?: string;
  level?: string;
  isStem?: boolean;
}

// ==================== BACKEND RESPONSE FORMATS ====================

// Legacy backend format (from existing database)
interface BackendInstitute {
  ID: string;
  Name: string;
  CreatedAt: string;
  UpdatedAt: string;
}

interface BackendFieldOfStudy {
  ID: string;
  Name: string;
  CreatedAt: string;
  UpdatedAt: string;
}

// New backend format (from Hipo API + enhanced database)
interface EnhancedBackendInstitute {
  ID?: string;
  Name: string;
  Country?: string;
  Domain?: string;
  WebPages?: string[];
  Source?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

interface EnhancedBackendFieldOfStudy {
  ID?: string;
  Code: string;
  Title: string;
  Name?: string; // Legacy field from old database
  Category?: string;
  Level?: string;
  IsStem?: boolean;
  CreatedAt?: string;
  UpdatedAt?: string;
}

// ==================== ADAPTER FUNCTIONS ====================

/**
 * Convert Institute to AutocompleteItem format
 */
function instituteToAutocompleteItem(institute: Institute): AutocompleteItem {
  return {
    id: institute.id || institute.name, // Use name as fallback ID for new items
    name: institute.name,
  };
}

/**
 * Convert FieldOfStudy to AutocompleteItem format
 * Handles both new format (Title) and legacy format (Name)
 */
function fieldOfStudyToAutocompleteItem(field: FieldOfStudy): AutocompleteItem {
  return {
    id: field.id || field.code, // Use code as fallback ID for new items
    name: field.title || field.name || "", // Fallback: Title -> Name -> empty
  };
}

// ==================== SERVICES ====================

export const instituteService = {
  /**
   * Search institutes from backend (with Hipo API + Redis caching)
   * Returns AutocompleteItem[] for compatibility with Autocomplete component
   * @param query - Search term (min 2 characters)
   * @param country - Optional country filter
   */
  search: async (
    query: string,
    country?: string
  ): Promise<AutocompleteItem[]> => {
    if (!query || query.length < 2) {
      return [];
    }

    try {
      const sanitized = sanitizeSearchQuery(query);
      const params: { q: string; country?: string } = { q: sanitized };

      if (country) {
        params.country = country;
      }

      const response = await api.get<{ results: EnhancedBackendInstitute[] }>(
        "/institutes/search",
        { params }
      );

      // Map backend format to frontend interface, then to AutocompleteItem
      const institutes: Institute[] = (response.data.results || []).map(
        (item) => ({
          id: item.ID,
          name: item.Name,
          country: item.Country,
          domain: item.Domain,
          webPages: item.WebPages,
          source: item.Source,
        })
      );

      return institutes.map(instituteToAutocompleteItem);
    } catch (error) {
      console.error("Failed to search institutes:", error);
      return [];
    }
  },

  /**
   * Search institutes with full data (not for autocomplete)
   * Use this when you need the complete Institute object with all fields
   */
  searchFull: async (query: string, country?: string): Promise<Institute[]> => {
    if (!query || query.length < 2) {
      return [];
    }

    try {
      const sanitized = sanitizeSearchQuery(query);
      const params: { q: string; country?: string } = { q: sanitized };

      if (country) {
        params.country = country;
      }

      const response = await api.get<{ results: EnhancedBackendInstitute[] }>(
        "/institutes/search",
        { params }
      );

      return (response.data.results || []).map((item) => ({
        id: item.ID,
        name: item.Name,
        country: item.Country,
        domain: item.Domain,
        webPages: item.WebPages,
        source: item.Source,
      }));
    } catch (error) {
      console.error("Failed to search institutes:", error);
      return [];
    }
  },

  /**
   * Get institute by ID (existing endpoint)
   */
  getById: async (id: string): Promise<Institute> => {
    const response = await api.get<BackendInstitute>(`/institutes/${id}`);
    return {
      id: response.data.ID,
      name: response.data.Name,
    };
  },
};

export const fieldOfStudyService = {
  /**
   * Search fields of study (CIP codes) from backend
   * Returns AutocompleteItem[] for compatibility with Autocomplete component
   * @param query - Search term (min 2 characters)
   * @param level - Optional level filter ("all" | "broad" | "detailed")
   */
  search: async (
    query: string,
    level: "all" | "broad" | "detailed" = "all"
  ): Promise<AutocompleteItem[]> => {
    if (!query || query.length < 2) {
      return [];
    }

    try {
      const sanitized = sanitizeSearchQuery(query);
      const response = await api.get<{
        results: EnhancedBackendFieldOfStudy[];
      }>("/field-of-studies/search", {
        params: { q: sanitized, level },
      });

      // Map backend format to frontend interface, then to AutocompleteItem
      const fields: FieldOfStudy[] = (response.data.results || []).map(
        (item) => ({
          id: item.ID,
          code: item.Code,
          title: item.Title,
          name: item.Name, // Legacy field
          category: item.Category,
          level: item.Level,
          isStem: item.IsStem,
        })
      );

      return fields.map(fieldOfStudyToAutocompleteItem);
    } catch (error) {
      console.error("Failed to search fields of study:", error);
      return [];
    }
  },

  /**
   * Search fields of study with full data (not for autocomplete)
   * Use this when you need the complete FieldOfStudy object with all fields
   */
  searchFull: async (
    query: string,
    level: "all" | "broad" | "detailed" = "all"
  ): Promise<FieldOfStudy[]> => {
    if (!query || query.length < 2) {
      return [];
    }

    try {
      const sanitized = sanitizeSearchQuery(query);
      const response = await api.get<{
        results: EnhancedBackendFieldOfStudy[];
      }>("/fields-of-study/search", {
        params: { q: sanitized, level },
      });

      return (response.data.results || []).map((item) => ({
        id: item.ID,
        code: item.Code,
        title: item.Title,
        name: item.Name, // Legacy field
        category: item.Category,
        level: item.Level,
        isStem: item.IsStem,
      }));
    } catch (error) {
      console.error("Failed to search fields of study:", error);
      return [];
    }
  },

  /**
   * Get field of study by ID (existing endpoint)
   */
  getById: async (id: string): Promise<FieldOfStudy> => {
    const response = await api.get<BackendFieldOfStudy>(
      `/field-of-studies/${id}`
    );
    return {
      id: response.data.ID,
      code: "", // Legacy data may not have code
      title: response.data.Name, // Legacy: Name is used as title
      name: response.data.Name, // Also populate name for consistency
    };
  },
};
