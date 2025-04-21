<script lang="ts">
    import type { LayoutProps } from '../admin2/$types';

    import Navbar from '$lib/components/navbar.svelte'
    let { data, children }: LayoutProps = $props();
</script>

<div class="">
    <Navbar/>
    <main>
        {@render children()}
    </main>
</div>

