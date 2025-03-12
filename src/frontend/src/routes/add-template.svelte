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
    IconPlus,
  } from "@tabler/icons-svelte";
  import Header from "../components/header.svelte";
  import * as system from "../../wailsjs/go/main/App";
  import { onMount } from "svelte";
  import { goto } from "@mateothegreat/svelte5-router";

  let containers = $state([]);
  let templates = $state([]);

  let volumes = $state([
    {
      source: "",
      destination: "/app",
    },
  ]);

  let name = $state("");
  let baseTemplate = $state("");
  let template = $state("");
  let isBusy = $state(false);

  onMount(async () => {
    templates = await system.GetTemplates(true);
    console.log("templates", templates);
  });

  const createTemplate = async () => {
    if (isBusy || !name.length || !baseTemplate.length || !template.length) {
      console.log("missing required fields");
      return;
    }

    isBusy = true;
    console.log("createTemplate", { name, baseTemplate, template });
    const error = await system.CreateTemplate(name, baseTemplate, template);

    if (error.length) {
      console.log("failed to create template", error);
    }

    isBusy = false;
    goto("/templates");
  };
</script>

<Header title="Add Template">
  <button
    onclick={createTemplate}
    disabled={isBusy}
    class="btn px-3 btn-success disabled:opacity-80">Create Template</button
  >
</Header>

<div class="flex flex-col gap-5">
  <div class="flex gap-5">
    <fieldset class="fieldset w-full">
      <legend class="fieldset-legend">Name</legend>
      <input
        bind:value={name}
        type="text"
        class="input w-full"
        placeholder="example-container"
      />
      <p class="fieldset-label">Required</p>
    </fieldset>
    <fieldset class="fieldset w-full">
      <legend class="fieldset-legend w-full">Base Template</legend>
      <select class="select w-full" bind:value={baseTemplate}>
        <option disabled value="">Pick a template</option>
        {#each templates as template}
          <option value={template.Id}>{template.Labels.name}</option>
        {/each}
      </select>
      <span class="fieldset-label">Required</span>
    </fieldset>
  </div>
  <fieldset class="fieldset">
    <legend class="fieldset-legend">Template</legend>
    <textarea
      bind:value={template}
      class="textarea h-24 w-full"
      placeholder="apt-get -y update && apt-get -y install xterm"
    ></textarea>
    <div class="fieldset-label">Required</div>
  </fieldset>
</div>
