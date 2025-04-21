<script lang="ts">
    import Link from "./link.svelte";

    let { children, onclick = () => {}, href = "" } = $props();
</script>

<Link
    {href}
    {onclick}
    class="inline-block p-1 rounded-sm cursor-pointer bg-primary-500 text-white font-bold hover:bg-transparent border-2 border-primary-500 hover:text-primary-500"
>
    {@render children?.()}
</Link>
