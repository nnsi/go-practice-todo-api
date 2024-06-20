const config = {
  API_URL: import.meta.env.VITE_API_URL || `http://${location.hostname}:8080`,
  WS_URL: import.meta.env.VITE_WS_URL || `ws://${location.hostname}:8080`,
};

export default config;
