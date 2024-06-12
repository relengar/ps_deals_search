import React from 'react';

export default function Button({
    text,
    className,
    children,
}: {
    text?: string;
    children?: React.ReactNode;
    className?: string;
}) {
    const classes = className ?? '';
    return (
        <button
            type="submit"
            className={`text-white bg-gray-800 hover:bg-gray-900 focus:outline-none focus:ring-4 focus:ring-gray-300 font-medium rounded-full text-sm dark:bg-gray-800 dark:hover:bg-gray-700 dark:focus:ring-gray-700 dark:border-gray-700 ${classes}`}
        >
            {text && text}
            {children && children}
        </button>
    );
}
