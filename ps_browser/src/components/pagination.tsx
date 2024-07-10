type PaginationProps = {
    total: number;
    limit: number;
    page: number;
    onChange: (page: number) => void;
};

export default function Pagination({
    total,
    limit,
    page,
    onChange,
}: PaginationProps) {
    const totalPages = Math.round(total / limit);
    const pageNumbers = getPageNumbers(page, totalPages);

    const getItemClasses = (pageNum: number) => {
        const selectedClasses = 'text-white bg-gray-700';
        const unselectedClasses = 'text-gray-400 bg-gray-800';
        if (pageNum === page) {
            return selectedClasses;
        }
        return unselectedClasses;
    };

    return (
        <ul className="flex items-center justify-center -space-x-px h-10 text-base">
            {page > 0 && (
                <>
                    <li>
                        <a
                            onClick={() => onChange(0)}
                            href="#"
                            className="flex items-center justify-center px-4 h-10 ms-0 leading-tight bg-gray-800 border-gray-700 text-gray-400 hover:bg-gray-700 hover:text-white"
                        >
                            <span className="sr-only">First</span>
                            <svg
                                className="w-7 h-7"
                                aria-hidden="true"
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                fill="none"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke="currentColor"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="m17 16-4-4 4-4m-6 8-4-4 4-4"
                                />
                            </svg>
                        </a>
                    </li>
                    <li>
                        <a
                            onClick={() => onChange(page - 1)}
                            href="#"
                            className="flex items-center justify-center px-4 h-10 ms-0 leading-tight bg-gray-800 border-gray-700 text-gray-400 hover:bg-gray-700 hover:text-white"
                        >
                            <span className="sr-only">Previous</span>
                            <svg
                                className="w-3 h-3"
                                aria-hidden="true"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 6 10"
                            >
                                <path
                                    stroke="currentColor"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="2"
                                    d="M5 1 1 5l4 4"
                                />
                            </svg>
                        </a>
                    </li>
                </>
            )}
            {pageNumbers.map((displayNum, num) => (
                // @ts-ignore
                <li key={num}>
                    <a
                        onClick={() => onChange(num)}
                        href="#"
                        className={`flex items-center justify-center px-4 h-10 leading-tight hover:text-gray-700 border-gray-700 hover:bg-gray-700 hover:text-white ${getItemClasses(
                            num
                        )}`}
                    >
                        {displayNum}
                    </a>
                </li>
            ))}
            {page < totalPages - 1 && (
                <>
                    <li>
                        <a
                            onClick={() => onChange(page + 1)}
                            href="#"
                            className="flex items-center justify-center px-4 h-10 ms-0 leading-tight bg-gray-800 border-gray-700 text-gray-400 hover:bg-gray-700 hover:text-white"
                        >
                            <span className="sr-only">Next</span>
                            <svg
                                className="w-7 h-7"
                                aria-hidden="true"
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                fill="none"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke="currentColor"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="m10 16 4-4-4-4"
                                />
                            </svg>
                        </a>
                    </li>
                    <li>
                        <a
                            onClick={() => onChange(totalPages - 1)}
                            href="#"
                            className="flex items-center justify-center px-4 h-10 ms-0 leading-tight bg-gray-800 border-gray-700 text-gray-400 hover:bg-gray-700 hover:text-white"
                        >
                            <span className="sr-only">Last</span>
                            <svg
                                className="w-7 h-7"
                                aria-hidden="true"
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                fill="none"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke="currentColor"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="m7 16 4-4-4-4m6 8 4-4-4-4"
                                />
                            </svg>
                        </a>
                    </li>
                </>
            )}
        </ul>
    );
}

function getPageNumbers(page: number, total: number) {
    const pagesToDisplay = 10;

    const start = Math.max(page - Math.floor(pagesToDisplay / 2), 0);
    const end = Math.min(page + pagesToDisplay - (page - start), total);

    const pageNumbers = [];
    for (let pageNum = start; pageNum < end; pageNum++) {
        pageNumbers.push(pageNum + 1);
    }

    return pageNumbers;
}
