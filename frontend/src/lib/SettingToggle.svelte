<script>
  import { createEventDispatcher } from 'svelte';
  
  export let label = '';
  export let description = '';
  export let checked = false;

  const dispatch = createEventDispatcher();

  function handleChange() {
    dispatch('change', !checked);
  }
</script>

<label class="flex items-center justify-between p-4 bg-[var(--bg-card)] rounded-lg cursor-pointer hover:bg-[var(--bg-hover)] transition-colors border border-[var(--border)]">
  <div>
    <p class="font-medium text-[var(--text-primary)]">{label}</p>
    {#if description}
      <p class="text-sm text-[var(--text-secondary)]">{description}</p>
    {/if}
  </div>
  
  <button 
    type="button"
    role="switch"
    aria-checked={checked}
    class="relative w-11 h-6 rounded-full transition-colors {checked ? 'bg-[var(--primary)]' : 'bg-[var(--bg-hover)]'}"
    on:click={handleChange}
  >
    <span 
      class="absolute top-1 left-1 w-4 h-4 bg-white rounded-full transition-transform shadow-sm {checked ? 'translate-x-5' : 'translate-x-0'}"
    ></span>
  </button>
</label>
