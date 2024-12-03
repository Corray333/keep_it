export interface Category {
    key: string;
    name: string;
    label: string;
    parentID: string;
    children: Category[];
    childrenCategories: Category[];
}