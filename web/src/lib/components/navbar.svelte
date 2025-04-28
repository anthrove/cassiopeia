<script lang="ts">
    import { useApiFetch } from "$lib/composables/apifetch";
    const api = useApiFetch();
    import { onMount } from "svelte";
    let navItems = [
        {
            name: "Home",
            link: "/",
            tenantRequired: false,
        },
        {
            name: "User Management",
            tenantRequired: true,
            link: "/users",
            subItems: [
                { name: "Groups", link: "/groups" },
                { name: "Users", link: "/users" },
            ],
        },
        {
            name: "Identity",
            tenantRequired: true,
            link: "/identity",
            subItems: [
                { name: "Applications", link: "/applications" },
                { name: "Certificates", link: "/certificates" },
                { name: "Provider", link: "/provider" },
            ],
        },
        {
            name: "Authorization",
            tenantRequired: true,
            link: "/authorization",
            subItems: [
                { name: "Permissions", link: "/permissions" },
                //{ name: "Models", link: "/models" },
                //{ name: "Adapters", link: "/adapters" },
                { name: "Enforcers", link: "/enforcers" },
            ],
        },
    ];

    import { readAll as readAllTenants } from "$lib/logic/tenant.svelte";
    let tenants: Tenant[] = $state([]);
    async function loadData() {
        const apitenants = await readAllTenants();
        if (!apitenants.map((t) => t.id).includes($context)) {
            context.set("");
        }
        tenants = <Tenant[]>apitenants || [
            {
                id: "",
                display_name: "-/-",
            },
        ];
    }
    onMount(loadData);

    import { context } from "$lib/stores/context.svelte";

    async function setCurentTenant(id: string) {
        context.set(id);
    }

    import Link from "./link.svelte";

    import { create as createCookieContext } from "$lib/stores/cookiecontext.svelte";
    let darkmode = createCookieContext("darkmode");
</script>

<nav class="shadow dark:bg-zinc-900">
    <nav class="p-4 container mx-auto flex gap-4 items-center">
        <b>
            <Link
                class="hover:before:content-['OwO'] before:text-rose-600 before:mr-1.5"
                href="/">WhoDis</Link
            >
        </b>
        <ul class="flex space-x-4 h-full z-10">
            {#each navItems as item}
                <li class="relative group h-full flex items-center">
                    {#if item.subItems}
                        <Link
                            href={item.link}
                            class="h-full flex items-center px-4"
                            disabled={item.tenantRequired && !$context}
                            >{item.name}</Link
                        >
                        {#if !item.tenantRequired || $context}
                            <ul
                                class="absolute left-4 top-4 pt-4 min-w-full mt-2 rounded shadow-lg hidden group-hover:block transition duration-300 ease-in-out transform opacity-0 group-hover:opacity-100 group-hover:translate-y-0"
                            >
                                {#each item.subItems as subItem}
                                    <li
                                        class="hover:bg-zinc-50 border-zinc-200 dark:border-none dark:hover:bg-zinc-700 border-b-[0.5px] bg-white dark:bg-zinc-900 "
                                    >
                                        <Link
                                            href={subItem.link}
                                            class="py-2 px-4 inline-block"
                                            >{subItem.name}</Link
                                        >
                                    </li>
                                {/each}
                            </ul>
                        {/if}
                    {:else}
                        <Link
                            href={item.link}
                            class="hover:text-gray-400 h-full flex items-center px-4"
                            >{item.name}</Link
                        >
                    {/if}
                </li>
            {/each}
        </ul>
        <div class="ml-auto flex gap-4 items-center">
            <li class="relative group h-full flex items-center">
                <Link
                    href="/tenant"
                    class="cursor-pointer h-full flex items-center px-4 italic"
                    >{tenants.find((t) => t.id == $context)?.display_name ||
                        "No tenant"}</Link
                >
                <ul
                    class="absolute right-0 top-4 pt-4 min-w-full mt-2 rounded shadow-lg hidden group-hover:block transition duration-300 ease-in-out transform opacity-0 group-hover:opacity-100 group-hover:translate-y-0"
                >
                    {#each tenants as tenant}
                        <li
                            class="hover:bg-zinc-50 border-zinc-200 border-b-[0.5px] bg-white text-black"
                        >
                            <Link
                                onclick={() => setCurentTenant(tenant.id)}
                                class="{$context == tenant.id
                                    ? 'text-primary-500 hover:text-primary-300'
                                    : 'hover:text-zinc-600'} py-2 px-4 inline-block cursor-pointer w-full"
                                >{tenant.display_name}</Link
                            >
                        </li>
                    {/each}
                </ul>
            </li>

            <div class="flex items-center">
                {#if $darkmode}
                    <button
                        onclick={() => darkmode.set(false)}
                        aria-label="toggle darkmode"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            ><path
                                fill="currentColor"
                                d="M9.37 5.51A7.4 7.4 0 0 0 9.1 7.5c0 4.08 3.32 7.4 7.4 7.4c.68 0 1.35-.09 1.99-.27A7.01 7.01 0 0 1 12 19c-3.86 0-7-3.14-7-7c0-2.93 1.81-5.45 4.37-6.49"
                                opacity="0.3"
                            /><path
                                fill="currentColor"
                                d="M9.37 5.51A7.4 7.4 0 0 0 9.1 7.5c0 4.08 3.32 7.4 7.4 7.4c.68 0 1.35-.09 1.99-.27A7.01 7.01 0 0 1 12 19c-3.86 0-7-3.14-7-7c0-2.93 1.81-5.45 4.37-6.49M12 3a9 9 0 1 0 9 9c0-.46-.04-.92-.1-1.36a5.39 5.39 0 0 1-4.4 2.26a5.403 5.403 0 0 1-3.14-9.8c-.44-.06-.9-.1-1.36-.1"
                            /></svg
                        >
                    </button>
                {:else}
                    <button
                        onclick={() => darkmode.set(true)}
                        aria-label="toggle light mode"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            ><path
                                fill="currentColor"
                                d="M12 7.5c-2.21 0-4 1.79-4 4s1.79 4 4 4s4-1.79 4-4s-1.79-4-4-4"
                                opacity="0.3"
                            /><path
                                fill="currentColor"
                                d="m5.34 6.25l1.42-1.41l-1.8-1.79l-1.41 1.41zM1 10.5h3v2H1zM11 .55h2V3.5h-2zm7.66 5.705l-1.41-1.407l1.79-1.79l1.406 1.41zM17.24 18.16l1.79 1.8l1.41-1.41l-1.8-1.79zM20 10.5h3v2h-3zm-8-5c-3.31 0-6 2.69-6 6s2.69 6 6 6s6-2.69 6-6s-2.69-6-6-6m0 10c-2.21 0-4-1.79-4-4s1.79-4 4-4s4 1.79 4 4s-1.79 4-4 4m-1 4h2v2.95h-2zm-7.45-.96l1.41 1.41l1.79-1.8l-1.41-1.41z"
                            /></svg
                        >
                    </button>
                {/if}
            </div>
        </div>
    </nav>
</nav>
