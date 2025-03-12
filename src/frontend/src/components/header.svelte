<script lang="ts">
  let { children, title }: { children: any, title: string } = $props()
</script>

<div class="border-b-2 border-white/10 p-4 flex justify-between items-center mb-8 mt-4">
  <div class="flex">
  <h2 class="text-2xl py-4">
    {title}
  </h2>
  </div>
      {@render children?.()}
</div>
