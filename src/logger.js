import log from "loglevel";

log.setDefaultLevel(import.meta.env.VITE_LOG_LEVEL || "info");

export default log;
