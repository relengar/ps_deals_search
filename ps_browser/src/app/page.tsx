import { redirect } from 'next/navigation';
import { Suspense } from 'react';
import LoadingGames from './loading';
import SearchGames from './search';
import Button from '@/components/button';
import SearchIcon from '@/components/searchIcon';

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
                    <div className="flex w-full mb-14 place-items-center justify-center gap-2">
                        <input
                            className="resize rounded-md center w-5/6 h-8"
                            placeholder="Search"
                            type="text"
                            name="term"
                        />
                        <Button
                            className="p-0"
                            children={
                                <SearchIcon className="h-8 w-8 p-2 m-0.5" />
                            }
                        />
                    </div>
                </form>
            </section>

            <Suspense fallback={<LoadingGames />}>
                <SearchGames term={term} />
            </Suspense>
        </section>
    );
}
