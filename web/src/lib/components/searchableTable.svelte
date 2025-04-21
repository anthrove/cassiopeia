<script lang="ts">

    type SearchableObject = {
        display_name: string;
        [key: string]: any;
    };

    let { rows, columns, row, header,empty } = $props();

    let search = $state("");

    let filteredItems = $derived(
        rows?.filter((item:SearchableObject) =>
            item.display_name?.toLowerCase().includes(search.toLowerCase()),
        ),
    );

    let columnsDerived = $derived(columns || Object.keys(rows[0] || {}));

    import Input from "./forms/input.svelte";
</script>

<div class="">
    {#if rows?.length}
    <div class="flex">
        <Input bind:value={search} placeholder="Search table"/>
    </div>
    {/if}
    <table class="w-full table-auto">
        <thead class="border-b-[0.5px] border-zinc-300">
            <tr class="">
                {#if header}
                    {@render header()}
                {:else}
                    {#each columnsDerived as column}
                        <th class="text-start">{column}</th>
                    {/each}
                {/if}
            </tr>
        </thead>
        <tbody>
            {#each filteredItems as item,i}
                <tr class="{i%2?'bg-zinc-100':''} border-b-[0.5px] border-zinc-300">
                    {@render row?.(item)}
                    {#if !row}
                        {#each columnsDerived as column}
                            <td class="py-4">{item[column]}</td>
                        {/each}
                    {/if}
                </tr>
            {/each}
        </tbody>
    </table>
    {#if !rows?.length}
        {@render empty?.()}
    {/if}
</div>