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
    <div class="relative flex items-center mb-1">
        <span class="bg-white px-2 mb-1 z-10 ml-8"
            >{label}{#if required}<span class="text-rose-600">*</span
                >{/if}</span
        >
        <div class="absolute w-full h-[0.5px] bg-zinc-200"></div>
    </div>
    <select disabled={readonly} {name} {placeholder} {required} class="{classList} w-full" bind:value={value}>
        {#each options as option}
            <option value="{option}">{labeler(option)}</option>
        {/each}
    </select>
</div>
