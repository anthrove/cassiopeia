<script lang="ts">
    import Link from "./link.svelte";

    type ButtonVariant = 'primary' | 'secondary' | 'danger'
    let {
        children,
        onclick = (e:PointerEvent) => {},
        href = "",
        class:classList = "",
        variant = 'primary' as ButtonVariant,
        type = 'button'
    } = $props();

    const VARIANTS = <{[key in ButtonVariant]:string}>{
        'primary':'bg-primary-500 text-white hover:bg-transparent border-[1px] border-primary-500 hover:text-primary-500',
        'danger':'bg-rose-600 text-white hover:font-bold hover:bg-transparent border-[1px] border-rose-600 hover:text-rose-600',
        'secondary':'border-[1px] border-zinc-500 text-black hover:bg-zinc-500 hover:text-white'
    }
</script>

<Link
    {href}
    {type}
    {onclick}
    class={`${classList} ${VARIANTS[variant]} p-1 px-2 cursor-pointer`}
>
    {@render children?.()}
</Link>