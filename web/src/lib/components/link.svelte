<script lang="ts">
    import { context } from "$lib/stores/context.svelte";

    const PREFIX = "/web";

    let {
        href,
        target = "",
        children,
        class: classList,
        disabled = false,
        onclick = () => {},
    } = $props();

    const computedHref = $derived.by(() => {
        if (!href) {
            return "";
        }
        if ($context) {
            return `${PREFIX}${href}?tenant=${$context}`;
        }
        return href;
    });

    let isButton = $derived(disabled || !href);
</script>

<svelte:element
    this={isButton ? "button" : "a"}
    onclick={() => {
        onclick();
    }}
    role={isButton ? "button" : "link"}
    href={computedHref}
    {target}
    class="{classList}{disabled
        ? ' cursor-not-allowed opacity-50'
        : ' cursor-pointer'} text-start">{@render children?.()}</svelte:element
>
