export default {
  // Consult https://github.com/sveltejs/svelte-preprocess
  // for more information about preprocessors
  vitePlugin: {
    dynamicCompileOptions: ({ filename }) => {
      if (filename.includes("node_modules/@tabler/icons-svelte")) {
        return { hydratable: true };
      }
      return {};
    },
  },
};
