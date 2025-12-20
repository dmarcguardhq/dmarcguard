import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import { useThemeStore, useSettingsStore } from "./stores";

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);

// Initialize stores before mounting
const themeStore = useThemeStore();
const settingsStore = useSettingsStore();

themeStore.init();
settingsStore.init();

app.mount("#app");
