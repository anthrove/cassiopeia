<script lang="ts">
    import type { LayoutProps } from './$types';

    let { data, children }: LayoutProps = $props();
    let navItems = [
        { name: 'Home', link: '/' },
        {
            name: 'User Management',
            subItems: [
                { name: 'Tenant', link: '/tenant' },
                { name: 'Groups', link: '/groups' },
                { name: 'User', link: '/user' }
            ]
        },
        {
            name: 'Identity',
            subItems: [
                { name: 'Applications', link: '/applications' },
                { name: 'Certificates', link: '/certificates' },
                { name: 'Provider', link: '/provider' }
            ]
        },
        {
            name: 'Authorization',
            subItems: [
                { name: 'Permissions', link: '/permissions' },
                { name: 'Models', link: '/models' },
                { name: 'Adapters', link: '/adapters' },
                { name: 'Enforcers', link: '/enforcers' }
            ]
        }
    ];
</script>


<div class="min-h-screen flex flex-col">
    <nav class="bg-gray-800 text-white p-4 flex justify-between items-center">
        <ul class="flex space-x-4 h-full">
            {#each navItems as item}
                <li class="relative group h-full flex items-center">
                    {#if item.subItems}
                        <div class="cursor-pointer h-full flex items-center px-4">{item.name}</div>
                        <ul class="absolute left-0 mt-2 bg-gray-700 text-white rounded shadow-lg hidden group-hover:block transition duration-300 ease-in-out transform opacity-0 group-hover:opacity-100 group-hover:translate-y-0">
                            {#each item.subItems as subItem}
                                <li class="py-1 px-4 hover:bg-gray-600">
                                    <a href={subItem.link} class="hover:text-gray-400">{subItem.name}</a>
                                </li>
                            {/each}
                        </ul>
                    {:else}
                        <a href={item.link} class="hover:text-gray-400 h-full flex items-center px-4">{item.name}</a>
                    {/if}
                </li>
            {/each}
        </ul>
    </nav>
    <main class="flex-grow p-4">
        {@render children()}
    </main>
</div>
