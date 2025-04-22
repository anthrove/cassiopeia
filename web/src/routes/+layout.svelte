<script lang="ts">
	import Navbar from '$lib/components/navbar.svelte'
	import type { LayoutProps } from './$types';

    export const prerender = true;
	import '../app.css';
	
	let { children } = $props();
</script>

<div class="">
    <Navbar/>
    <main class="container p-4 mx-auto">
        {@render children()}
    </main>
</div>