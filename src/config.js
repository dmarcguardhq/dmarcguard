var BASE_API_URL = import.meta.env.VITE_BASE_API_URL || "/api";
const isDev = import.meta.env.DEV;

if (isDev) {
  BASE_API_URL = "http://localhost:8080/api";
}

export var BASE_API_URL;
