<script lang="ts">
    import DateInput from "./DateInput.svelte";

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
</script>

<label class="flex flex-col {classList}">
    {#if label}
        <div class="relative flex items-center mb-1">
            <span class="px-2 z-10 ml-8 dark:bg-zinc-800 bg-white"
                >{label}{#if required}<span class="text-rose-600">*</span
                    >{/if}</span
            >
            <div
                class="absolute w-full h-[1px] bg-zinc-200 dark:bg-zinc-600"
            ></div>
        </div>
    {/if}
    {#if ["text", "number"].includes(inputType)}
        <input
            {readonly}
            type={inputType}
            {required}
            {name}
            placeholder={placeholder || label}
            bind:value
            class="bg-inherit {readonly
                ? 'cursor-not-allowed text-zinc-400'
                : ''} p-2 border-zinc-200"
        />
    {:else if ["checkbox"].includes(inputType)}
        <div class="flex items-center gap-2">
            <div class="w-[40px]">
                <label class="inline-flex items-center cursor-pointer">
                    <input
                        class="sr-only peer"
                        bind:checked={value}
                        type="checkbox"
                        {name}
                        {readonly}
                        {required}
                    />
                    <div
                        class="relative w-11 h-6 bg-zinc-300 dark:bg-zinc-700 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-primary-500 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"
                    ></div>
                </label>
            </div>
            
        </div>
    {:else if ["date", "datetime-local"].includes(inputType)}
        <DateInput
            {readonly}
            type={inputType}
            {required}
            {name}
            placeholder={placeholder || label}
            bind:value
            class="{readonly
                ? 'cursor-not-allowed text-zinc-400'
                : ''} p-2 border-zinc-200"
        />
    {:else}
        <div class="text-lg text-rose-600 italic">
            Input type not yet implemented :(
        </div>
    {/if}
</label>
