#!/bin/bash -e

rm -rf EAST || true
git clone https://github.com/argman/EAST
cd EAST
# get east_icdar2015_resnet_v1_50_rbox.zip
rm east_icdar2015_resnet_v1_50_rbox.zip || true
wget -O east_icdar2015_resnet_v1_50_rbox.zip "https://drive.google.com/uc?id=0B3APw5BZJ67ETHNPaU9xUkVoV0U&export=download&confirm=t&uuid=3b82612e-5f89-4639-938b-3d5ba0010dac"
cp ../freeze_east_model.py .

unzip ./east_icdar2015_resnet_v1_50_rbox.zip
docker run -u $(id -u):$(id -g) -v ${PWD}/:/EAST:rw -w /EAST tensorflow/tensorflow:1.15.5 python3 freeze_east_model.py
docker run -u $(id -u):$(id -g) -v ${PWD}/:/EAST:rw openvino/ubuntu20_dev:2022.3.0 mo \
--framework=tf --input_shape=[1,1024,1920,3] --input=input_images --output=feature_fusion/Conv_7/Sigmoid,feature_fusion/concat_3 \
--input_model /EAST/model.pb --output_dir /EAST/IR/1/

# download text-recognition model
curl -L --create-dir https://storage.openvinotoolkit.org/repositories/open_model_zoo/2022.1/models_bin/2/text-recognition-0014/FP32/text-recognition-0014.bin -o text-recognition/1/model.bin https://storage.openvinotoolkit.org/repositories/open_model_zoo/2022.1/models_bin/2/text-recognition-0014/FP32/text-recognition-0014.xml -o text-recognition/1/model.xml
chmod -R 755 text-recognition/

mkdir -p ../../config/OCR/east_fp32
mkdir -p ../../config/OCR/text-recognition

rm -rf ../../config/OCR/east_fp32
rm -rf ../../config/OCR/text-recognition
cp -R ./IR ../../config/OCR/east_fp32
cp -R ./text-recognition ../../config/OCR/

echo "build east ocr share library..."
cd ../../src/custom_nodes
make NODES=east_ocr
rm -rf ../../config/OCR/lib
cp -R ./lib/ubuntu ../../config/OCR/lib

cd ../../ocr

