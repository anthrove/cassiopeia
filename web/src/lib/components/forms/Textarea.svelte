<script lang="ts">
    import { onMount } from "svelte";
    let textarea = $state<HTMLTextAreaElement>(<any>null);
    let {
        value = $bindable(<any>""),
        label = "",
        type: inputType = "text",
        placeholder = "",
        required = false,
        name = "",
        class: classList = "",
        readonly = false,
    } = $props();

    function resize() {
        if (!textarea) {
            return;
        }
        textarea.style.height = "auto"; // Reset height
        textarea.style.height = `${textarea.scrollHeight}px`; // Set to scrollHeight
    }

    function handleInput(event) {
        value = event.target.value;
        resize();
    }

    onMount(resize);
</script>

{#if label}
    <div class="relative flex items-center mb-1">
        <span class="px-2 z-10 ml-8 dark:bg-zinc-800 bg-white"
            >{label}{#if required}<span class="text-rose-600">*</span
                >{/if}</span
        >
        <div class="absolute w-full h-[1px] bg-zinc-200 dark:bg-zinc-600"></div>
    </div>
{/if}
<textarea
    bind:this={textarea}
    bind:value
    oninput={handleInput}
    rows="3"
    class="{classList} bg-inherit {readonly
        ? 'cursor-not-allowed text-zinc-400'
        : ''} p-2 border-zinc-200 w-full"
></textarea>
