<script lang="ts">
    import Crud from "$lib/pages/CRUD.svelte";

    import * as apiControllerUntyped from "$lib/logic/tenant.svelte";
    const apiController = apiControllerUntyped as unknown as CRUDAPI;

    const OBJECT_TYPE = "Tenant";
    const TABLE_COLUMNS = ["Tennant", "Password Type"];

    import apidescriptor from '../../data/apidescriptor.json'
    const build_createForm_descriptor = {
        display_name: String,
        password_type: apidescriptor.tenant.password_types,
    }

    import { readAll as readAllCertificates } from "$lib/logic/certificate.svelte";
    async function build_editForm_descriptor(){
        const certificates = await readAllCertificates();

        return <FormDescriptor>{
            id: {
                type: "string",
                label: "Tenant ID",
                readonly: true,
            },
            display_name: {
                type: "string",
                label: "Display name",
                required: true,
            },
            password_type: {
                type: "select:single",
                label: "Password type",
                required: true,

                options: apidescriptor.tenant.password_types,
            },
            signing_certificate_id: {
                type: "select:single",
                labeler: (id: string) =>
                    `${certificates.find((c: Certificate) => c.id == id)?.display_name} (${id})`,
                label: "Signing key",
                required: true,

                options: certificates.map((c) => c.id),
            },

            created_at: {
                type: "string",
                label: "Created at",
                readonly: true,
            },
            updated_at: {
                type: "string",
                label: "Last update",
                readonly: true,
            },
        };
    }
</script>

<Crud
    {apiController}
    {OBJECT_TYPE}
    {TABLE_COLUMNS}
    {build_createForm_descriptor}
    {build_editForm_descriptor}
>
    {#snippet tennant(item: User)}
        {item.display_name}
    {/snippet}
</Crud>
