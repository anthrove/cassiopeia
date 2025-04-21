<script lang="ts">
    let {
        value = $bindable(""),
        label = "",
        type: inputType = "text",
        placeholder = "",
        required = false,
        name = "",
        class:classList=''
    } = $props();
</script>

{#if ["text"].includes(inputType)}
    <label class="flex flex-col {classList}">
        {#if label}
        <div class="relative flex items-center mb-1">
            <span class="bg-white px-2 mb-1 z-10 ml-8"
                >{label}{#if required}<span class="text-rose-600">*</span
                    >{/if}</span
            >
            <div class="absolute w-full h-[0.5px] bg-zinc-200"></div>
        </div>
        {/if}
        <input
            type="text"
            {required}
            {name}
            placeholder={placeholder || label}
            bind:value
            class="p-2 border-zinc-200"
        />
    </label>
{:else}
    <div class="text-lg text-rose-600 italic">Not yet implemented :(</div>
{/if}
