import FileSaver from 'file-saver';

export const SaveFile = (path: string, contents: string) => {
    return FileSaver.saveAs(contents, path);
}
