augroup project
  autocmd!
  autocmd BufRead,BufNewFile *.h,*.c set filetype=c
  autocmd BufNewFile *.h HeaderguardAdd
augroup END

let &path.="include,"

nnoremap <C-M> :make!<CR>
nnoremap <Leader>h :CocCommand clangd.switchSourceHeader<CR>
