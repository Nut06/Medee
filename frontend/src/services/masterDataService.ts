import { api } from "@/lib/api";
import { sanitizeSearchQuery } from "@/utils/validation";

export interface Institute {
  id: string;
  name: string;
}

export interface FieldOfStudy {
  id: string;
  name: string;
}

// Backend response format (from Go)
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

export const instituteService = {
  search: async (query: string): Promise<Institute[]> => {
    if (query.length < 2) {
      return [];
    }

    const sanitized = sanitizeSearchQuery(query);
    const response = await api.get<{ results: BackendInstitute[] }>(
      "/institutes/search",
      {
        params: { q: sanitized },
      }
    );

    // Map backend format to frontend interface
    return (response.data.results || []).map((item) => ({
      id: item.ID,
      name: item.Name,
    }));
  },

  getById: async (id: string): Promise<Institute> => {
    const response = await api.get<BackendInstitute>(`/institutes/${id}`);
    return {
      id: response.data.ID,
      name: response.data.Name,
    };
  },
};

export const fieldOfStudyService = {
  search: async (query: string): Promise<FieldOfStudy[]> => {
    if (query.length < 2) {
      return [];
    }

    const sanitized = sanitizeSearchQuery(query);
    const response = await api.get<{ results: BackendFieldOfStudy[] }>(
      "/field-of-studies/search",
      {
        params: { q: sanitized },
      }
    );

    // Map backend format to frontend interface
    return (response.data.results || []).map((item) => ({
      id: item.ID,
      name: item.Name,
    }));
  },

  getById: async (id: string): Promise<FieldOfStudy> => {
    const response = await api.get<BackendFieldOfStudy>(
      `/field-of-studies/${id}`
    );
    return {
      id: response.data.ID,
      name: response.data.Name,
    };
  },
};
