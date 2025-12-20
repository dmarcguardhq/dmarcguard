import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";

const STORAGE_KEY = "theme";
const VALID_THEMES = ["light", "dark", "system"];

/**
 * Theme store for managing application theme with localStorage persistence.
 * Supports light, dark, and system (auto-detect) modes.
 */
export const useThemeStore = defineStore("theme", () => {
  // State
  const theme = ref("system");
  const systemPrefersDark = ref(false);

  // Media query for system preference detection
  let mediaQuery = null;

  /**
   * Resolved theme value (light or dark) based on current setting
   */
  const resolvedTheme = computed(() => {
    if (theme.value === "system") {
      return systemPrefersDark.value ? "dark" : "light";
    }
    return theme.value;
  });

  /**
   * Check if dark mode is currently active
   */
  const isDark = computed(() => resolvedTheme.value === "dark");

  /**
   * Initialize theme from localStorage and setup system preference listener
   */
  function init() {
    // Setup system preference detection
    if (typeof window !== "undefined" && window.matchMedia) {
      mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
      systemPrefersDark.value = mediaQuery.matches;

      // Listen for system preference changes
      mediaQuery.addEventListener("change", (e) => {
        systemPrefersDark.value = e.matches;
      });
    }

    // Load saved theme from localStorage
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved && VALID_THEMES.includes(saved)) {
      theme.value = saved;
    }

    // Apply theme to document
    applyTheme();
  }

  /**
   * Apply the resolved theme to the document
   */
  function applyTheme() {
    document.documentElement.setAttribute("data-theme", resolvedTheme.value);
  }

  /**
   * Set the theme and persist to localStorage
   * @param {string} newTheme - One of 'light', 'dark', or 'system'
   */
  function setTheme(newTheme) {
    if (!VALID_THEMES.includes(newTheme)) {
      console.warn(
        `Invalid theme: ${newTheme}. Must be one of: ${VALID_THEMES.join(", ")}`,
      );
      return;
    }
    theme.value = newTheme;
    localStorage.setItem(STORAGE_KEY, newTheme);
  }

  /**
   * Cycle through themes: light -> dark -> system -> light
   */
  function cycleTheme() {
    const currentIndex = VALID_THEMES.indexOf(theme.value);
    const nextIndex = (currentIndex + 1) % VALID_THEMES.length;
    setTheme(VALID_THEMES[nextIndex]);
  }

  // Watch for theme changes and apply to document
  watch(resolvedTheme, () => {
    applyTheme();
  });

  return {
    // State
    theme,
    resolvedTheme,
    isDark,
    // Actions
    init,
    setTheme,
    cycleTheme,
  };
});
