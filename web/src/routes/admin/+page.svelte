<script lang="ts">
  import { onMount } from 'svelte';
  import { writable } from 'svelte/store';

  interface Tenant {
    id: string;
    display_name: string;
    created_at: string;
    updated_at: string;
    password_type: string;
    signing_key_id: string;
  }

  let tenants: Tenant[] = [];
  let selectedTenant: Tenant | null = null;
  let loading = writable(true);
  let error = writable<string | null>(null);

  onMount(async () => {
    loading.set(true);
    error.set(null);
    try {
      const res = await fetch('/api/v1/tenant');
      if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
      const json = await res.json();
      if (json.error) throw new Error(json.error);
      tenants = json.data ?? [];
    } catch (e) {
      error.set(e instanceof Error ? e.message : String(e));
    } finally {
      loading.set(false);
    }
  });

  function selectTenant(tenant: Tenant) {
    selectedTenant = tenant;
  }
</script>

<div class="p-4">
  <h2 class="text-xl font-bold mb-4">Tenants</h2>

  {#if $loading}
    <p>Loading tenants...</p>
  {:else if $error}
    <p class="text-red-600">Error: {$error}</p>
  {:else if tenants.length === 0}
    <p>No tenants found.</p>
  {:else}
    <table class="min-w-full border-collapse border border-gray-300">
      <thead>
      <tr class="bg-gray-200">
        <th class="border border-gray-300 px-4 py-2 text-left">ID</th>
        <th class="border border-gray-300 px-4 py-2 text-left">Display Name</th>
        <th class="border border-gray-300 px-4 py-2 text-left">Created At</th>
        <th class="border border-gray-300 px-4 py-2 text-left">Updated At</th>
        <th class="border border-gray-300 px-4 py-2 text-left">Password Type</th>
        <th class="border border-gray-300 px-4 py-2 text-left">Signing Key ID</th>
      </tr>
      </thead>
      <tbody>
      {#each tenants as tenant}
        <tr
                class="cursor-pointer hover:bg-gray-100 {selectedTenant?.id === tenant.id ? 'bg-blue-100' : ''}"
                on:click={() => selectTenant(tenant)}
        >
          <td class="border border-gray-300 px-4 py-2">{tenant.id}</td>
          <td class="border border-gray-300 px-4 py-2">{tenant.display_name}</td>
          <td class="border border-gray-300 px-4 py-2">{new Date(tenant.created_at).toLocaleString()}</td>
          <td class="border border-gray-300 px-4 py-2">{new Date(tenant.updated_at).toLocaleString()}</td>
          <td class="border border-gray-300 px-4 py-2">{tenant.password_type}</td>
          <td class="border border-gray-300 px-4 py-2">{tenant.signing_key_id}</td>
        </tr>
      {/each}
      </tbody>
    </table>

    {#if selectedTenant}
      <div class="mt-6 p-4 border rounded bg-gray-50 space-y-1">
        <h3 class="font-semibold">Selected Tenant Details</h3>
        <p><strong>ID:</strong> {selectedTenant.id}</p>
        <p><strong>Display Name:</strong> {selectedTenant.display_name}</p>
        <p><strong>Created At:</strong> {new Date(selectedTenant.created_at).toLocaleString()}</p>
        <p><strong>Updated At:</strong> {new Date(selectedTenant.updated_at).toLocaleString()}</p>
        <p><strong>Password Type:</strong> {selectedTenant.password_type}</p>
        <p><strong>Signing Key ID:</strong> {selectedTenant.signing_key_id}</p>
      </div>
    {/if}
  {/if}
</div>

<style>
  table {
    border-collapse: collapse;
  }
  th, td {
    border: 1px solid #d1d5db; /* Tailwind gray-300 */
    padding: 0.5rem 1rem;
  }
  tr:hover {
    background-color: #f3f4f6; /* Tailwind gray-100 */
  }
  tr.selected {
    background-color: #bfdbfe; /* Tailwind blue-200 */
  }
</style>
