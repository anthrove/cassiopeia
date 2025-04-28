<script lang="ts">
	import Navbar from '$lib/components/navbar.svelte'
	import type { LayoutProps } from './$types';

    export const prerender = true;
	import '../app.css';
	
	let { children } = $props();

    import { create as createCookieContext } from "$lib/stores/cookiecontext.svelte";
    let darkmode = createCookieContext("darkmode");
    let darkmodeClass = $derived($darkmode?'dark':'light');

    import { classList, setBody } from "$lib/actions/classList";
    $effect(()=>{
        setBody(darkmodeClass)
    })
</script>

<svelte:body/>

<div class="">
    <Navbar/>
    <main class="container p-4 mx-auto">
        {@render children()}
    </main>
</div>