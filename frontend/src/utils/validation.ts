export const sanitizeSearchQuery = (query: string): string => {
  // Remove dangerous characters
  return query
    .trim()
    .replace(/[<>\"']/g, '') // Remove HTML/SQL injection chars
    .substring(0, 100);      // Limit length
};

export const validateUUID = (id: string): boolean => {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
  return uuidRegex.test(id);
};