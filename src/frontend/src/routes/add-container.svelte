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
  let name = $state("");
  let template = $state("");
  let startupCommand = $state("");
  let isBusy = $state(false);

  let volumes = $state([]);

  onMount(async () => {
    containers = await system.Containers();
    templates = await system.GetTemplates(true);
    console.log("containers", containers);
  });

  const createContainer = async () => {
    if (isBusy || !name || !template || !startupCommand) {
      console.error("missing required fields");
      return;
    }

    isBusy = true;

    const filteredVolumes = volumes
      .filter((volume) => volume.source?.length && volume.destination?.length)
      .map((volume) => {
        return `${volume.source}:${volume.destination}`;
      });

    console.log("filteredVolumes", filteredVolumes);

    const response = await system.CreateContainer(
      name,
      template,
      startupCommand,
      filteredVolumes,
    );

    isBusy = false;
    console.log("response", response);
    goto("/");
  };

  const addVolume = () => {
    volumes = [...volumes, { source: "", destination: "/app" }];
  };

  const deleteVolume = (index) => {
    volumes = volumes.filter((_, i) => i !== index);
  };
</script>

<Header title="Add Container">
  <button
    onclick={createContainer}
    class="btn px-3 btn-success disabled:opacity-80"
  >
    Create Container
  </button>
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
      <legend class="fieldset-legend">Template</legend>
      <select class="select w-full" bind:value={template}>
        <option disabled selected value="">Pick a template</option>
        {#each templates as template}
          <option value={template.Id}>{template.Labels.name}</option>
        {/each}
      </select>
      <span class="fieldset-label">Required</span>
    </fieldset>
  </div>
  <fieldset class="fieldset">
    <legend class="fieldset-legend w-full">Startup Command</legend>
    <input
      bind:value={startupCommand}
      type="text"
      class="input w-full"
      placeholder="xterm"
    />
    <div class="fieldset-label">Required</div>
  </fieldset>
  {#each volumes as volume, index}
    <div class="flex gap-5">
      <fieldset class="fieldset w-full">
        <legend class="fieldset-legend w-full">Volume Source</legend>
        <input
          type="text"
          class="input w-full"
          onchange={(e) => (volumes[index].source = e.target.value)}
          value={volume.source}
          placeholder="/root/example"
        />
      </fieldset>
      <fieldset class="fieldset w-full flex">
        <legend class="fieldset-legend w-full">Volume Destination</legend>
        <input
          type="text"
          value={volume.destination}
          onchange={(e) => (volumes[index].destination = e.target.value)}
          class="input w-full"
          placeholder="/app"
        />
        <div class="flex items-end">
          <button class="btn btn-danger" onclick={() => deleteVolume(index)}>
            <IconTrash size={24} />
          </button>
        </div>
      </fieldset>
    </div>
  {/each}
  <div class="flex gap-5 justify-between">
    <button class="btn btn-primary" onclick={() => addVolume()}
      >Add Volume</button
    >
  </div>
</div>
