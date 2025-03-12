<script lang="ts">
  import type { RouterInstance } from "@mateothegreat/svelte5-router";
  import { Router, goto } from "@mateothegreat/svelte5-router";
  import { menuItems, routes } from "./constants";
  import { onMount } from "svelte";
  import * as system from "../wailsjs/go/main/App";
  import * as runtime from "../wailsjs/runtime/runtime.js";

  let { route } = $props();
  let instance = $state<RouterInstance>();
  let logs = $state([]);
  let path = $state("/");

  onMount(async () => {
    runtime.EventsOn("syslog", (data: any) => {
      console.log("received syslog", data);
      logs = [data, ...logs];
    });

    runtime.EventsOn("openLog", () =>
      document.getElementById("log_modal").showModal(),
    );

    runtime.EventsOn("closeLog", () => {
      document.getElementById("log_modal").close();
      logs = [];
    });

    await system.Init();
  });
</script>

<main class="flex">
  <div
    id="sidebar"
    class="w-[25vw] h-[100vh] bg-secondary/10 flex flex-col items-center justify-between"
  >
    <div class="w-full flex flex-col items-center">
      <h1 class="text-4xl py-10 w-full text-center">AppJar</h1>
      {#each menuItems as item}
        <a
          href={item.path}
          onclick={(e) => {
            e.preventDefault();
            path = item.path;
            goto(item.path);
          }}
          class:opacity-100={path === item.path}
          class={`w-full text-center hover:bg-secondary/20 block p-4 text-white opacity-70 text-2xl active:opacity-100 transform-gpu transition-bg ${
            path === item.path ? "bg-secondary/20" : ""
          }`}
        >
          <item.icon class="w-10 h-10 inline-block mr-4" />
          {item.name}
        </a>
      {/each}
    </div>
    <div class="mb-4">
      <p class="text-md">AppJar [v0.1]</p>
    </div>
  </div>
  <div id="content" class="w-[75vw] h-[100vh] min-w-[880px] mx-8">
    <Router bind:instance {routes} />
  </div>
  <dialog id="log_modal" class="modal">
    <div class="modal-box">
      <h3 class="text-lg font-bold mb-4">Logs</h3>
      <div class="max-h-[50vh] overflow-auto">
        {#each logs as log}
          <p>{log}</p>
        {/each}
      </div>
    </div>
  </dialog>
</main>
