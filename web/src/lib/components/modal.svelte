<script lang="ts">
    let { open = $bindable(false), children, onclose=()=>{} } = $props();

    function closeSelf(event:Event) {
        if (event.target === event.currentTarget) {
            onclose()
            open = false
        }
    }
    function close(){
        onclose()
        open = false
    }
</script>

{#if open}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <div role="dialog" tabindex="0" class="absolute overflow-y-auto h-screen flex items-start inset-0 bg-black/40 z-20 text-start" onclick={closeSelf}>
        <div class="mt-16 p-4 flex-1">
            <div class="container mx-auto p-4 bg-white rounded shadow block relative">
                <button aria-label="close" class="absolute top-2 right-2 p-2 cursor-pointer hover:text-primary-500" onclick={close}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M8.4 17L7 15.6l3.6-3.6L7 8.425l1.4-1.4l3.6 3.6l3.575-3.6l1.4 1.4l-3.6 3.575l3.6 3.6l-1.4 1.4L12 13.4z"/></svg>
                </button>
                {@render children?.()}
            </div>
        </div>
    </div>
{/if}
