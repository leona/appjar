import { IconBox, IconEdit, IconSettings } from "@tabler/icons-svelte";
import { type RouteConfig } from "@mateothegreat/svelte5-router";
import Home from "./routes/home.svelte";
import Settings from "./routes/settings.svelte";
import AddContainer from "./routes/add-container.svelte";
import Templates from "./routes/templates.svelte";
import AddTemplate from "./routes/add-template.svelte";

export const menuItems = [
  {
    name: "Containers",
    path: "/",
    icon: IconBox,
  },
  {
    name: "Templates",
    path: "/templates",
    icon: IconEdit,
  },
  {
    name: "Settings",
    path: "/settings",
    icon: IconSettings,
  },
];

export const routes: RouteConfig[] = [
  {
    component: Home,
  },
  {
    component: Settings,
    path: "/settings",
  },
  {
    component: AddContainer,
    path: "/containers/add",
  },
  {
    component: AddTemplate,
    path: "/add-template",
  },
  {
    component: Templates,
    path: "/templates",
  },
];
