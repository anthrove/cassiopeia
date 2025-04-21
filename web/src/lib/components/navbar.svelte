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

    let tenants: Tenant[] = [];
    onMount(async () => {
        const apitenants = (await api<Tenant[]>("v1/tenant")) || [];
        if (!apitenants.map((t) => t.id).includes($context)) {
            context.set("");
        }
        tenants = <Tenant[]>apitenants||[
            {
                id: "",
                display_name: "-/-",
            },
        ];
    });

    import { context } from "$lib/stores/context.svelte";

    async function setCurentTenant(id: string) {
        context.set(id);
    }

    import Link from "./link.svelte";
    import Button from "./Button.svelte";
</script>

<nav class="p-4 bg-white text-black shadow flex gap-4 items-center">
    <b>
        <Link
            class="hover:before:content-['OwO'] before:text-rose-600 before:mr-1.5"
            href="/">WhoDis</Link
        >
    </b>
    <ul class="flex space-x-4 h-full">
        {#each navItems as item}
            <li class="relative group h-full flex items-center">
                {#if item.subItems}
                    <Link
                        href={item.link}
                        class="h-full flex items-center px-4"
                        disabled={item.tenantRequired&&!$context}
                        >{item.name}</Link
                    >
                    {#if !item.tenantRequired || $context}
                    <ul
                        class="absolute left-4 top-4 pt-4 min-w-full mt-2 rounded shadow-lg hidden group-hover:block transition duration-300 ease-in-out transform opacity-0 group-hover:opacity-100 group-hover:translate-y-0"
                    >
                        {#each item.subItems as subItem}
                            <li
                                class="hover:bg-zinc-50 border-zinc-200 border-b-[0.5px] bg-white text-black"
                            >
                                <Link
                                    href={subItem.link}
                                    class="hover:text-zinc-600 py-2 px-4 inline-block"
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
    <div class="ml-auto">
        <li class="relative group h-full flex items-center">
            <Link
                disabled={!$context}
                href="/tenant/{$context}"
                class="cursor-pointer h-full flex items-center px-4 italic"
                >{tenants.find(t=>t.id == $context)?.display_name || "No tenant"}</Link
            >
            <ul
                class="absolute right-0 top-4 pt-4 min-w-full mt-2 rounded shadow-lg hidden group-hover:block transition duration-300 ease-in-out transform opacity-0 group-hover:opacity-100 group-hover:translate-y-0"
            >
                {#each tenants as tenant}
                    <li
                        class="hover:bg-zinc-50 border-zinc-200 border-b-[0.5px] bg-white text-black"
                    >
                        <Link
                            href="{$context == tenant.id?`/tenant/${$context}`:''}"
                            onclick={() => setCurentTenant(tenant.id)}
                            class="{$context == tenant.id?'text-primary-500 hover:text-primary-300':'hover:text-zinc-600'} py-2 px-4 inline-block cursor-pointer w-full"
                            >{tenant.display_name}</Link
                        >
                    </li>
                {/each}
            </ul>
        </li>
    </div>
</nav>
