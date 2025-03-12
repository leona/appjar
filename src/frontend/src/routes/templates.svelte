<script lang="ts">
  import { IconTrash } from "@tabler/icons-svelte";
  import Header from "../components/header.svelte";
  import * as system from "../../wailsjs/go/main/App";
  import { onMount, onDestroy } from "svelte";
  import { goto } from "@mateothegreat/svelte5-router";

  let templates = $state([]);
  let isBusy = $state(false);
  let interval;

  onMount(async () => {
    templates = await system.GetTemplates(true);
    console.log("templates", templates);

    interval = setInterval(async () => {
      if (!isBusy) {
        templates = await system.GetTemplates(true);
      }
    }, 3000);
  });

  const deleteTemplate = async (id: string) => {
    isBusy = true;
    const error = await system.DeleteTemplate(id);
    templates = await system.GetTemplates(true);
    isBusy = false;
  };

  onDestroy(() => {
    clearInterval(interval);
  });
</script>

<Header title="Templates">
  <button onclick={() => goto("/add-template")} class="btn px-3 btn-success">
    Add
  </button>
</Header>

<div>
  {#if templates.length === 0}
    <p class="text-white/70 text-center py-8">No templates found</p>
  {/if}
  {#each templates as template}
    <div
      class="bg-secondary/30 flex justify-between px-6 py-2 my-4 rounded-sm items-center"
    >
      <div class="flex items-center gap-8">
        <div class="flex flex-col">
          <p class="text-sm text-white/70">
            {template.Labels.type}
          </p>
          <h3 class="text-xl">{template.Labels.name}</h3>
        </div>
      </div>
      <div class="flex gap-3">
        <button
          onclick={() =>
            document
              .getElementById("confirm_delete_modal" + template.Id)
              .showModal()}
          disabled={isBusy}
          class="btn btn-error text-white"
        >
          <IconTrash class="w-8 h-8" />
        </button>
      </div>
    </div>
    <dialog id={"confirm_delete_modal" + template.Id} class="modal">
      <div class="modal-box">
        <h3 class="text-lg font-bold">Confirm</h3>
        <p class="py-4">
          Are you sure you want to delete the template {template.Labels.name}?
        </p>
        <div class="flex justify-end gap-2">
          <button
            class="btn btn-error text-white"
            disabled={isBusy}
            onclick={async () => {
              await deleteTemplate(template.Id);
              document
                .getElementById("confirm_delete_modal" + template.Id)
                .close();
            }}>Delete</button
          >
          <button
            class="btn btn-secondary text-white"
            onclick={() =>
              document
                .getElementById("confirm_delete_modal" + template.Id)
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
