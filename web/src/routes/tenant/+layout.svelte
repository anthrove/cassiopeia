<script lang="ts">
    import Nocontext from "$lib/components/errors/nocontext.svelte";
    import { context } from "$lib/stores/context.svelte";

    let {children} = $props()
</script>

{#if $context}
    {@render children?.()}
{:else}
    <Nocontext/>
{/if}