import { redirect } from 'next/navigation';
import { Suspense } from 'react';
import LoadingGames from './loading';
import SearchGames from './search';
import Button from '@/components/button';

type QueryParams = {
    term?: string;
};

export default async function Search({
    searchParams,
}: {
    searchParams: QueryParams;
}) {
    const { term } = searchParams;

    const goToTerm = async (form: FormData) => {
        'use server';
        const term = form.get('term')?.toString();
        if (term?.length === 0) {
            return;
        }
        redirect(`?term=${term}`);
    };

    return (
        <section className="container mx-auto flex-col content-center">
            <h2 className="text-center w-full text-3xl leading-loose p-10">
                Search for deals
            </h2>
            <section>
                <form action={goToTerm}>
                    <div className="flex w-full mb-14 justify-evenly place-items-center">
                        <input
                            className="resize rounded-md center w-5/6 h-8"
                            placeholder="Search"
                            type="text"
                            name="term"
                        />

                        <Button text={'Search'} />
                    </div>
                </form>
            </section>

            <Suspense fallback={<LoadingGames />}>
                <SearchGames term={term} />
            </Suspense>
        </section>
    );
}
