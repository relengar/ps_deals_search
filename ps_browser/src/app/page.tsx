import { searchGames } from '@/lib/queries/searchGamesQuery';
import { redirect } from 'next/navigation';

type QueryParams = {
    term?: string;
};

export default async function Search({
    searchParams,
}: {
    searchParams: QueryParams;
}) {
    const { term } = searchParams;

    let data: any[] | null = null;
    if (term) {
        const { games } = await searchGames({ term });
        data = games;
    }

    const goToTerm = async (form: FormData) => {
        'use server';
        const term = form.get('term')?.toString();
        if (term?.length === 0) {
            return;
        }
        redirect(`?term=${term}`);
    };

    return (
        <section>
            <h3>Search</h3>
            <div>
                <form action={goToTerm}>
                    <input placeholder="Search" type="text" name="term" />
                    <button>Search</button>
                    <div>
                        Searched for <span>{term}</span>
                    </div>
                    <div>
                        {data && data.map((game) => <div>{game.name}</div>)}
                    </div>
                </form>
            </div>
        </section>
    );
}
