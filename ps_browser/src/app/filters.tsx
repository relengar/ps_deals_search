import { FiltersIcon } from '@/components/icons/filtersIcon';
import { SearchGameParams } from '@/lib/queries/searchGames';
import type { Platform } from '@/lib/repositories/games';
import { ChangeEvent, useState } from 'react';

type Props = {
    term: string;
    onChange: (filters: Partial<SearchGameParams>) => void;
} & SearchGameParams;

export default function Filters(props: Props) {
    const { onChange } = props;
    const [maxPrice, setMaxPrice] = useState(props.maxPrice);
    const setPrice = (evt: ChangeEvent<HTMLInputElement>) => {
        const price = Number(evt.target.value);
        setMaxPrice(price);
        onChange({ maxPrice: price });
    };

    const [useSemantic, setUseSemantic] = useState(props.useSemantic);
    const toggleSemantic = () => {
        const use = !useSemantic;
        setUseSemantic(use);
        onChange({ useSemantic: use });
    };

    const [visible, setVisible] = useState(false);
    const toggleVisible = () => setVisible(!visible);

    const [platforms, setPlatforms] = useState<Platform[]>(
        props.platforms ?? []
    );
    const onPlatformsChange = (evt: ChangeEvent<HTMLInputElement>) => {
        const selected = evt.target.value as Platform;
        const newPlatforms = toSelected(selected, platforms);

        // cannot have no platforms selected
        if (newPlatforms.length === 0) {
            return;
        }

        setPlatforms(newPlatforms);
        onChange({ platforms: newPlatforms });
    };

    return (
        <>
            <div className="flex w-full mb-6 place-items-center justify-center gap-2">
                <span className=" w-5/6"></span>
                <a href="#" onClick={toggleVisible}>
                    <FiltersIcon on={visible} className="h-10 w-10 p-2 m-0.5" />
                </a>
            </div>
            <div
                className={`${
                    visible ? 'visible' : 'invisible'
                } flex w-full mb-14 place-items-center justify-center gap-2`}
            >
                <div className="flex items-center ">
                    <label
                        htmlFor="default-range"
                        className="block text-sm font-medium dark:text-white"
                    >
                        Max price
                    </label>
                    <div className="relative ml-2">
                        <input
                            id="default-range"
                            type="range"
                            min="0"
                            max="500"
                            value={maxPrice}
                            onChange={setPrice}
                            className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
                        />
                        <span className="text-sm absolute start-1/2 -translate-x-1/2 rtl:translate-x-1/2 -bottom-6">
                            {maxPrice}
                        </span>
                    </div>
                </div>

                <div className="flex items-center">
                    <label className="inline-flex items-center cursor-pointer">
                        <input
                            type="checkbox"
                            value=""
                            className="sr-only peer"
                            disabled={props.term.length === 0}
                            checked={useSemantic}
                            onChange={toggleSemantic}
                        />
                        <div className="relative w-11 h-6 bg-gray-200 rounded-full peer peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                        <span className="ms-3 text-sm font-medium">
                            Search in description
                        </span>
                    </label>
                </div>

                <div className="flex items-center">
                    <div className="flex items-center ps-4">
                        <input
                            id="bordered-checkbox-1"
                            type="checkbox"
                            value="PS4"
                            name="PS4"
                            className="w-4 h-4 bg-gray-100 border-gray-600 rounded focus:ring-blue-600 ring-offset-gray-800 focus:ring-2"
                            onChange={onPlatformsChange}
                            checked={platforms.includes('PS4')}
                        />
                        <label
                            htmlFor="PS4"
                            className="w-full py-4 ms-2 text-sm font-medium text-gray-300"
                        >
                            PS 4
                        </label>
                    </div>
                    <div className="flex items-center ps-4">
                        <input
                            id="bordered-checkbox-1"
                            type="checkbox"
                            value="PS5"
                            name="PS5"
                            className="w-4 h-4 bg-gray-100 border-gray-600 rounded focus:ring-blue-600 ring-offset-gray-800 focus:ring-2"
                            onChange={onPlatformsChange}
                            checked={platforms.includes('PS5')}
                        />
                        <label
                            htmlFor="PS5"
                            className="w-full py-4 ms-2 text-sm font-medium text-gray-300"
                        >
                            PS 5
                        </label>
                    </div>
                </div>
            </div>
        </>
    );
}

function toSelected(selected: Platform, current: Platform[]): Platform[] {
    if (current.includes(selected)) {
        return current.filter((p) => p !== selected);
    }

    return [...current, selected];
}
