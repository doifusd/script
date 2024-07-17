import pandas as pd
import os

def convert_xls_to_xlsx(src_folder, dest_folder):
    """ 
    转换 src_folder 中的所有.xls 文件为.xlsx 格式，并保存在 dest_folder 中 
    """
    # 遍历源文件夹中的所有文件
    for filename in os.listdir(src_folder):
        if filename.endswith('.xls'):
            # 构建完整的文件路径
            file_path = os.path.join(src_folder, filename) 
            # 读取.xls 文件
            df = pd.read_excel(file_path) 
            # 构建新的文件名和路径
            new_filename = filename.replace('.xls', '.xlsx')
            new_file_path = os.path.join(dest_folder, new_filename) 
            # 将数据写入.xlsx 文件
            df.to_excel(new_file_path, index=False) 

# 指定源文件夹（存放 xls 文件的文件夹）和目标文件夹（转换后的 xlsx 文件存放处）
src_folder = '/Users/sky/Documents/script/excel/data/data1' 
dest_folder = '/Users/sky/Documents/script/merge/data' 
convert_xls_to_xlsx(src_folder, dest_folder)
