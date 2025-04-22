<script lang="ts">
    import Link from "../link.svelte";

</script>

<div class="flex gap-8 items-center max-w-max mx-auto">
    <div class="text-rose-600">
        <svg
            xmlns="http://www.w3.org/2000/svg"
            width="200"
            height="200"
            viewBox="0 0 24 24"
            ><path
                fill="currentColor"
                d="M12 17q.425 0 .713-.288T13 16t-.288-.712T12 15t-.712.288T11 16t.288.713T12 17m-1-4h2V7h-2zm1 9q-2.075 0-3.9-.788t-3.175-2.137T2.788 15.9T2 12t.788-3.9t2.137-3.175T8.1 2.788T12 2t3.9.788t3.175 2.137T21.213 8.1T22 12t-.788 3.9t-2.137 3.175t-3.175 2.138T12 22"
            /></svg
        >
    </div>
    <div class="prose">
        <h1 class="text-rose-600">Context missing</h1>
        <p>
            Please select what tenant you want to interact with<br>
            <Link href="/">Go Home</Link>
        </p>
    </div>
</div>
