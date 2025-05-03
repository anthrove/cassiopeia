<script lang="ts">
    let {
        value = $bindable(""),
        options = $bindable([]),
        label = <string|LabelerFunction>"",
        labeler= <LabelerFunction>(id:string)=>id,
        type: inputType = "text",
        placeholder = "",
        required = false,
        name = "",
        class:classList='',
        readonly=false
    } = $props();
</script>

<div class="">
    {#if label}
    <div class="relative flex items-center mb-1">
        <span class="bg-white dark:bg-zinc-800 px-2 mb-1 z-10 ml-8"
            >{label}{#if required}<span class="text-rose-600">*</span
                >{/if}</span
        >
        <div class="absolute w-full h-[1px] bg-zinc-200 dark:bg-zinc-600"></div>
    </div>
    {/if}
    <select disabled={readonly} {name} {placeholder} {required} class="{classList} border-zinc-200 bg-inherit w-full" bind:value={value}>
        {#each options as option}
            <option value="{option}">{labeler(option)}</option>
        {/each}
    </select>
</div>
