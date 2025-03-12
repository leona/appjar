<script lang="ts">
  import {
    IconSettings,
    IconPlayerPlay,
    IconPlayerStop,
    IconTrash,
    IconAppWindow,
    IconDeviceDesktopPin,
    IconWindowMaximize,
    IconEdit,
  } from "@tabler/icons-svelte";
  import Header from "../components/header.svelte";
  import * as system from "../../wailsjs/go/main/App";
  import { onMount, onDestroy } from "svelte";
  import { goto } from "@mateothegreat/svelte5-router";

  let containers = $state([]);
  let templates = $state([]);
  let isBusy = $state(false);
  let interval;

  onMount(async () => {
    containers = await system.Containers();
    templates = await system.GetTemplates(true);
    console.log("containers", containers, "templates", templates);

    interval = setInterval(async () => {
      if (!isBusy) {
        containers = await system.Containers();
      }
    }, 3000);
  });

  const templateName = (id: string) => {
    const template = templates.find((t) => t.Id === id);
    return template?.Labels?.name;
  };

  const setContainerState = async (id: string, state: "start" | "stop") => {
    isBusy = true;
    const error = await system.SetContainerState(id, state);
    containers = await system.Containers();
    isBusy = false;

    if (error?.length) {
      console.error(error);
    }
  };

  const deleteContainer = async (id: string) => {
    isBusy = true;
    const error = await system.DeleteContainer(id);
    containers = await system.Containers();
    isBusy = false;
  };

  const connectContainer = async (container: any) => {
    isBusy = true;
    await system.XpraConnect(container.Labels.port, container.Labels.password);
    isBusy = false;
  };

  const openWeb = async (container: any) => {
    isBusy = true;
    await system.OpenLink(`http://localhost:${container.Labels.port}`);
    isBusy = false;
  };

  onDestroy(() => {
    clearInterval(interval);
  });
</script>

<Header title="Containers">
  <button onclick={() => goto("/containers/add")} class="btn px-3 btn-success">
    Add
  </button>
</Header>

<div>
  {#if containers.length === 0}
    <p class="text-white/70 text-center py-8">No containers found</p>
  {/if}
  {#each containers as container}
    <div
      class="bg-secondary/30 flex justify-between px-6 py-2 my-4 rounded-sm items-center"
    >
      <div class="flex items-center gap-8">
        <div class="flex flex-col">
          <p class="text-sm text-white/70">
            {templateName(container.Labels.template)}
          </p>
          <h3 class="text-xl">{container.Labels.name}</h3>
        </div>
        {#if container.State === "running"}
          <p class="bg-success/60 text-2xl py-1 px-4 rounded-sm">Running</p>
        {/if}
      </div>
      <div class="flex gap-3">
        {#if container.State === "running"}
          <div class="mr-4 flex gap-2">
            <button
              disabled={isBusy}
              onclick={() => openWeb(container)}
              class="btn btn-primary text-white"
            >
              <IconWindowMaximize class="w-8 h-8" />
            </button>
            <!--
            <button class="btn btn-primary text-white">
              <IconDeviceDesktopPin class="w-8 h-8" />
            </button>
            -->
            <button
              onclick={() => connectContainer(container)}
              disabled={isBusy}
              class="btn btn-primary text-white"
            >
              <IconAppWindow class="w-8 h-8" />
            </button>
          </div>
          <button
            disabled={isBusy}
            onclick={() => setContainerState(container.Id, "stop")}
            class="btn btn-accent text-white"
          >
            <IconPlayerStop class="w-8 h-8" />
          </button>
        {:else}
          <button
            class="btn btn-success text-white"
            disabled={isBusy}
            onclick={() => setContainerState(container.Id, "start")}
          >
            <IconPlayerPlay class="w-8 h-8" />
          </button>
        {/if}
        <button
          class="btn btn-error text-white"
          onclick={() =>
            document
              .getElementById("confirm_delete_modal" + container.Id)
              .showModal()}
        >
          <IconTrash class="w-8 h-8" />
        </button>
      </div>
    </div>
    <dialog id={"confirm_delete_modal" + container.Id} class="modal">
      <div class="modal-box">
        <h3 class="text-lg font-bold">Confirm</h3>
        <p class="py-4">
          Are you sure you want to delete the container {container.Labels.name}?
        </p>
        <div class="flex justify-end gap-2">
          <button
            class="btn btn-error text-white"
            disabled={isBusy}
            onclick={async () => {
              await deleteContainer(container.Id);
              document
                .getElementById("confirm_delete_modal" + container.Id)
                .close();
            }}>Delete</button
          >
          <button
            class="btn btn-secondary text-white"
            onclick={() =>
              document
                .getElementById("confirm_delete_modal" + container.Id)
                .close()}>Cancel</button
          >
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  {/each}
</div>
