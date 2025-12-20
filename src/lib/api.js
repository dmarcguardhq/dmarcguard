import ky from "ky";

/**
 * Configured ky instance for API requests.
 *
 * Features:
 * - Automatic JSON parsing
 * - Consistent error handling
 * - Retry logic for transient failures
 * - Configurable timeout
 */
const api = ky.create({
  prefixUrl: "/api",
  timeout: 30000,
  retry: {
    limit: 2,
    methods: ["get"],
    statusCodes: [408, 500, 502, 503, 504],
  },
  hooks: {
    beforeError: [
      (error) => {
        const { response } = error;
        if (response?.body) {
          error.message = `API Error: ${response.status} ${response.statusText}`;
        }
        return error;
      },
    ],
  },
});

/**
 * Statistics API
 */
export const getStatistics = () => api.get("statistics").json();

/**
 * Top sources API
 * @param {number} limit - Maximum number of sources to return
 */
export const getTopSources = (limit = 10) =>
  api.get("top-sources", { searchParams: { limit } }).json();

/**
 * Reports API
 * @param {Object} options - Query options
 * @param {number} options.limit - Maximum number of reports to return
 * @param {number} options.offset - Offset for pagination
 */
export const getReports = ({ limit = 20, offset = 0 } = {}) =>
  api.get("reports", { searchParams: { limit, offset } }).json();

/**
 * Get a single report by ID
 * @param {string|number} id - Report ID
 */
export const getReportById = (id) => api.get(`reports/${id}`).json();

export default api;
